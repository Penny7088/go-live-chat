package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/model"
)

var _ InterestsDao = (*interestsDao)(nil)

// InterestsDao defining the dao interface
type InterestsDao interface {
	Create(ctx context.Context, table *model.Interests) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Interests) error
	GetByID(ctx context.Context, id uint64) (*model.Interests, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Interests, int64, error)
	GetByLanguage(ctx context.Context, params string) ([]*model.InterestsTranslations, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.Interests, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.Interests, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.Interests, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Interests) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Interests) error
}

type interestsDao struct {
	db    *gorm.DB
	cache cache.InterestsCache // if nil, the cache is not used.
	sfg   *singleflight.Group  // if cache is nil, the sfg is not used.
}

// NewInterestsDao creating the dao interface
func NewInterestsDao(db *gorm.DB, xCache cache.InterestsCache) InterestsDao {
	if xCache == nil {
		return &interestsDao{db: db}
	}
	return &interestsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *interestsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *interestsDao) Create(ctx context.Context, table *model.Interests) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *interestsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Interests{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *interestsDao) UpdateByID(ctx context.Context, table *model.Interests) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *interestsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Interests) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.TagID != 0 {
		update["tag_id"] = table.TagID
	}
	if table.TagName != "" {
		update["tag_name"] = table.TagName
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *interestsDao) GetByID(ctx context.Context, id uint64) (*model.Interests, error) {
	// no cache
	if d.cache == nil {
		record := &model.Interests{}
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
			table := &model.Interests{}
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
			err = d.cache.Set(ctx, id, table, cache.InterestsExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Interests)
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
func (d *interestsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Interests, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Interests{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Interests{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

func (d *interestsDao) GetByLanguage(ctx context.Context, params string) ([]*model.InterestsTranslations, error) {
	var interests []*model.InterestsTranslations
	err := d.db.WithContext(ctx).Table("interests_translations as it").Select("it.tag_id, it.translated_name").Joins("JOIN interests as i ON it.tag_id = i.id").Where("it.language_code = ?", params).Find(&interests).Error
	if err != nil {
		return nil, err
	}
	return interests, nil
}

// DeleteByIDs delete records by batch id
func (d *interestsDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.Interests{}).Error
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
func (d *interestsDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.Interests, error) {
	queryStr, args, err := c.ConvertToGorm()
	if err != nil {
		return nil, err
	}

	table := &model.Interests{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs get records by batch id
func (d *interestsDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.Interests, error) {
	// no cache
	if d.cache == nil {
		var records []*model.Interests
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.Interests)
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
			var missedData []*model.Interests
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
			if err != nil {
				return nil, err
			}

			if len(missedData) > 0 {
				for _, data := range missedData {
					itemMap[data.ID] = data
				}
				err = d.cache.MultiSet(ctx, missedData, cache.InterestsExpireTime)
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
func (d *interestsDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.Interests, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.Interests{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *interestsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Interests) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *interestsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.Interests{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *interestsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Interests) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
