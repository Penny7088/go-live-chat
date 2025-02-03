package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"lingua_exchange/internal/types"

	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/model"
)

var _ GroupNoticeDao = (*groupNoticeDao)(nil)

// GroupNoticeDao defining the dao interface
type GroupNoticeDao interface {
	Create(ctx context.Context, table *model.GroupNotice) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.GroupNotice) error
	GetByID(ctx context.Context, id uint64) (*model.GroupNotice, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.GroupNotice, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.GroupNotice, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.GroupNotice, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.GroupNotice, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.GroupNotice) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.GroupNotice) error
	GetListAll(ctx context.Context, groupId int) ([]*types.SearchNoticeItem, error)
}

type groupNoticeDao struct {
	db    *gorm.DB
	cache cache.GroupNoticeCache // if nil, the cache is not used.
	sfg   *singleflight.Group    // if cache is nil, the sfg is not used.
}

// NewGroupNoticeDao creating the dao interface
func NewGroupNoticeDao(db *gorm.DB, xCache cache.GroupNoticeCache) GroupNoticeDao {
	if xCache == nil {
		return &groupNoticeDao{db: db}
	}
	return &groupNoticeDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *groupNoticeDao) GetListAll(ctx context.Context, groupId int) ([]*types.SearchNoticeItem, error) {
	fields := []string{
		"group_notice.id",
		"group_notice.creator_id",
		"group_notice.title",
		"group_notice.content",
		"group_notice.is_top",
		"group_notice.is_confirm",
		"group_notice.confirm_users",
		"group_notice.created_at",
		"group_notice.updated_at",
		"users.avatar",
		"users.profile_picture",
	}
	query := d.db.WithContext(ctx).Table("group_notice")
	query.Joins("left join users on users.id = group_notice.creator_id")
	query.Where("group_notice.group_id = ? and group_notice.is_delete = ?", groupId, 0)
	query.Order("group_notice.is_top desc")
	query.Order("group_notice.created_at desc")

	var items []*types.SearchNoticeItem
	if err := query.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}

func (d *groupNoticeDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *groupNoticeDao) Create(ctx context.Context, table *model.GroupNotice) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *groupNoticeDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.GroupNotice{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *groupNoticeDao) UpdateByID(ctx context.Context, table *model.GroupNotice) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *groupNoticeDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.GroupNotice) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.GroupID != 0 {
		update["group_id"] = table.GroupID
	}
	if table.CreatorID != 0 {
		update["creator_id"] = table.CreatorID
	}
	if table.Title != "" {
		update["title"] = table.Title
	}
	if table.Content != "" {
		update["content"] = table.Content
	}
	if table.ConfirmUsers != "" {
		update["confirm_users"] = table.ConfirmUsers
	}
	if table.IsDelete != 0 {
		update["is_delete"] = table.IsDelete
	}
	if table.IsTop != 0 {
		update["is_top"] = table.IsTop
	}
	if table.IsConfirm != 0 {
		update["is_confirm"] = table.IsConfirm
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *groupNoticeDao) GetByID(ctx context.Context, id uint64) (*model.GroupNotice, error) {
	// no cache
	if d.cache == nil {
		record := &model.GroupNotice{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache or database
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	if errors.Is(err, model.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { // nolint
			table := &model.GroupNotice{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				// if data is empty, set not found cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, model.ErrRecordNotFound) {
					err = d.cache.SetCacheWithNotFound(ctx, id)
					if err != nil {
						return nil, err
					}
					return nil, model.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			err = d.cache.Set(ctx, id, table, cache.GroupNoticeExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.GroupNotice)
		if !ok {
			return nil, model.ErrRecordNotFound
		}
		return table, nil
	} else if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, model.ErrRecordNotFound
	}

	// fail fast, if cache error return, don't request to db
	return nil, err
}

// GetByColumns get paging records by column information,
// Note: query performance degrades when table rows are very large because of the use of offset.
//
// params includes paging parameters and query parameters
// paging parameters (required):
//
//	page: page number, starting from 0
//	limit: lines per page
//	sort: sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
//
// query parameters (not required):
//
//	name: column name
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like, in
//	value: column value, if exp=in, multiple values are separated by commas
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: search for a male over 20 years of age
//
//	params = &query.Params{
//	    Page: 0,
//	    Limit: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *groupNoticeDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.GroupNotice, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.GroupNotice{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.GroupNotice{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs delete records by batch id
func (d *groupNoticeDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.GroupNotice{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a record by condition
// query conditions:
//
//	name: column name
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like, in
//	value: column value, if exp=in, multiple values are separated by commas
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: find a male aged 20
//
//	condition = &query.Conditions{
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *groupNoticeDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.GroupNotice, error) {
	queryStr, args, err := c.ConvertToGorm()
	if err != nil {
		return nil, err
	}

	table := &model.GroupNotice{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs get records by batch id
func (d *groupNoticeDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.GroupNotice, error) {
	// no cache
	if d.cache == nil {
		var records []*model.GroupNotice
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.GroupNotice)
		for _, record := range records {
			itemMap[record.ID] = record
		}
		return itemMap, nil
	}

	// get form cache or database
	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	var missedIDs []uint64
	for _, id := range ids {
		_, ok := itemMap[id]
		if !ok {
			missedIDs = append(missedIDs, id)
			continue
		}
	}

	// get missed data
	if len(missedIDs) > 0 {
		// find the id of an active placeholder, i.e. an id that does not exist in database
		var realMissedIDs []uint64
		for _, id := range missedIDs {
			_, err = d.cache.Get(ctx, id)
			if errors.Is(err, cacheBase.ErrPlaceholder) {
				continue
			}
			realMissedIDs = append(realMissedIDs, id)
		}

		if len(realMissedIDs) > 0 {
			var missedData []*model.GroupNotice
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
			if err != nil {
				return nil, err
			}

			if len(missedData) > 0 {
				for _, data := range missedData {
					itemMap[data.ID] = data
				}
				err = d.cache.MultiSet(ctx, missedData, cache.GroupNoticeExpireTime)
				if err != nil {
					return nil, err
				}
			} else {
				for _, id := range realMissedIDs {
					_ = d.cache.SetCacheWithNotFound(ctx, id)
				}
			}
		}
	}

	return itemMap, nil
}

// GetByLastID get paging records by last id and limit
func (d *groupNoticeDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.GroupNotice, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.GroupNotice{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *groupNoticeDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.GroupNotice) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *groupNoticeDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.GroupNotice{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *groupNoticeDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.GroupNotice) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
