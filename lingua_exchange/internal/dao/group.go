package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lingua_exchange/internal/types"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/model"
)

var _ GroupDao = (*groupDao)(nil)

// GroupDao defining the dao interface
type GroupDao interface {
	Create(ctx context.Context, table *model.Group) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Group) error
	GetByID(ctx context.Context, id uint64) (*model.Group, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Group, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.Group, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.Group, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.Group, error)
	GroupList(ctx context.Context, id int) ([]*types.GroupItem, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Group) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Group) error
	SearchOvertList(ctx context.Context, opt *types.SearchOvertListOpt) ([]*model.Group, error)
	FindAll(ctx context.Context, arg ...func(*gorm.DB)) ([]*model.Group, error)
}

type groupDao struct {
	db    *gorm.DB
	cache cache.GroupCache    // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

func (d *groupDao) SearchOvertList(ctx context.Context, opt *types.SearchOvertListOpt) ([]*model.Group, error) {
	return d.FindAll(ctx, func(db *gorm.DB) {
		if opt.Name != "" {
			db.Where("name like ?", "%"+opt.Name+"%")
		}
		db.Where("is_overt = ?", 1)
		db.Where("id NOT IN (?)", d.db.Select("group_id").Where("user_id = ? and is_quit= ?", opt.UserId, 0).Table("group_member"))
		db.Where("is_dismiss = 0").Order("created_at desc").Offset((opt.Page - 1) * opt.Size).Limit(opt.Size)
	})
}

func (d *groupDao) GroupList(ctx context.Context, id int) ([]*types.GroupItem, error) {
	tx := d.db.Table("group_member")
	tx.Select("`group`.id,`group`.name as group_name,`group`.avatar,`group`.profile,group_member.leader,`group`.creator_id")
	tx.Joins("left join `group` on `group`.id = group_member.group_id")
	tx.Where("group_member.user_id = ? and group_member.is_quit = ?", id, 0)
	tx.Order("group_member.created_at desc")
	items := make([]*types.GroupItem, 0)
	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}
	length := len(items)
	if length == 0 {
		return items, nil
	}

	ids := make([]int, 0, length)
	for i := range items {
		ids = append(ids, items[i].ID)
	}

	query := d.db.Table("talk_session")
	query.Select("receiver_id,is_disturb")
	query.Where("talk_type = ? and receiver_id in ?", 2, ids)
	list := make([]*types.TalkSessionPart, 0)
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	hash := make(map[int]*types.TalkSessionPart)
	for i := range list {
		hash[list[i].ReceiverID] = list[i]
	}

	for i := range items {
		if value, ok := hash[items[i].ID]; ok {
			items[i].IsDisturb = value.IsDisturb
		}
	}

	return items, nil
}

func (d *groupDao) FindAll(ctx context.Context, arg ...func(*gorm.DB)) ([]*model.Group, error) {
	db := d.db.Model(ctx)
	for _, fn := range arg {
		fn(db)
	}
	var items []*model.Group
	if err := db.Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// NewGroupDao creating the dao interface
func NewGroupDao(db *gorm.DB, xCache cache.GroupCache) GroupDao {
	if xCache == nil {
		return &groupDao{db: db}
	}
	return &groupDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *groupDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *groupDao) Create(ctx context.Context, table *model.Group) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *groupDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Group{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *groupDao) UpdateByID(ctx context.Context, table *model.Group) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, uint64(table.ID))

	return err
}

func (d *groupDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Group) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Type != 0 {
		update["type"] = table.Type
	}
	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Profile != "" {
		update["profile"] = table.Profile
	}
	if table.Avatar != "" {
		update["avatar"] = table.Avatar
	}
	if table.MaxNum != 0 {
		update["max_num"] = table.MaxNum
	}
	if table.IsOvert != 0 {
		update["is_overt"] = table.IsOvert
	}
	if table.IsMute != 0 {
		update["is_mute"] = table.IsMute
	}
	if table.IsDismiss != 0 {
		update["is_dismiss"] = table.IsDismiss
	}
	if table.CreatorID != 0 {
		update["creator_id"] = table.CreatorID
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *groupDao) GetByID(ctx context.Context, id uint64) (*model.Group, error) {
	// no cache
	if d.cache == nil {
		record := &model.Group{}
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
			table := &model.Group{}
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
			err = d.cache.Set(ctx, id, table, cache.GroupExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Group)
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
func (d *groupDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Group, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Group{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Group{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs delete records by batch id
func (d *groupDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.Group{}).Error
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
func (d *groupDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.Group, error) {
	queryStr, args, err := c.ConvertToGorm()
	if err != nil {
		return nil, err
	}

	table := &model.Group{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs get records by batch id
func (d *groupDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.Group, error) {
	// no cache
	if d.cache == nil {
		var records []*model.Group
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.Group)
		for _, record := range records {
			itemMap[uint64(record.ID)] = record
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
			var missedData []*model.Group
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
			if err != nil {
				return nil, err
			}

			if len(missedData) > 0 {
				for _, data := range missedData {
					itemMap[uint64(data.ID)] = data
				}
				err = d.cache.MultiSet(ctx, missedData, cache.GroupExpireTime)
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
func (d *groupDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.Group, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.Group{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *groupDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Group) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return uint64(table.ID), err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *groupDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.Group{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *groupDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Group) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, uint64(table.ID))

	return err
}
