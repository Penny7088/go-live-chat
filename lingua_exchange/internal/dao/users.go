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

var _ UsersDao = (*usersDao)(nil)

// UsersDao defining the dao interface
type UsersDao interface {
	Create(ctx context.Context, table *model.Users) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Users) error
	GetByID(ctx context.Context, id uint64) (*model.Users, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Users, int64, error)
	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.Users, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.Users, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.Users, error)
	GetByEmail(ctx context.Context, params string) (*model.Users, error)
	GetByEmailTx(ctx context.Context, tx *gorm.DB, params *model.Users) (*model.Users, error)
	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Users) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Users) error
}

type usersDao struct {
	db    *gorm.DB
	cache cache.UsersCache    // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewUsersDao creating the dao interface
func NewUsersDao(db *gorm.DB, xCache cache.UsersCache) UsersDao {
	if xCache == nil {
		return &usersDao{db: db}
	}
	return &usersDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *usersDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *usersDao) Create(ctx context.Context, table *model.Users) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *usersDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Users{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *usersDao) UpdateByID(ctx context.Context, table *model.Users) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *usersDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Users) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Email != "" {
		update["email"] = table.Email
	}
	if table.Username != "" {
		update["username"] = table.Username
	}
	if table.PasswordHash != "" {
		update["password_hash"] = table.PasswordHash
	}
	if table.ProfilePicture != "" {
		update["profile_picture"] = table.ProfilePicture
	}
	if table.NativeLanguageID != 0 {
		update["native_language_id"] = table.NativeLanguageID
	}
	if table.Age != 0 {
		update["age"] = table.Age
	}
	if table.Gender != "" {
		update["gender"] = table.Gender
	}
	if table.CountryID != 0 {
		update["country_id"] = table.CountryID
	}
	if table.RegistrationDate.IsZero() == false {
		update["registration_date"] = table.RegistrationDate
	}
	if table.LastLogin.IsZero() == false {
		update["last_login"] = table.LastLogin
	}
	if table.Status != "" {
		update["status"] = table.Status
	}
	if table.EmailVerified != 0 {
		update["email_verified"] = table.EmailVerified
	}
	if table.VerificationToken != "" {
		update["verification_token"] = table.VerificationToken
	}
	if table.TokenExpiration.IsZero() == false {
		update["token_expiration"] = table.TokenExpiration
	}
	if table.BirthDate.IsZero() == false {
		update["birth_date"] = table.BirthDate
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *usersDao) GetByID(ctx context.Context, id uint64) (*model.Users, error) {
	// no cache
	if d.cache == nil {
		record := &model.Users{}
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
			table := &model.Users{}
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
			err = d.cache.Set(ctx, id, table, cache.UsersExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Users)
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
func (d *usersDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Users, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Users{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Users{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs delete records by batch id
func (d *usersDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.Users{}).Error
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
func (d *usersDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.Users, error) {
	queryStr, args, err := c.ConvertToGorm()
	if err != nil {
		return nil, err
	}

	table := &model.Users{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs get records by batch id
func (d *usersDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.Users, error) {
	// no cache
	if d.cache == nil {
		var records []*model.Users
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.Users)
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
			var missedData []*model.Users
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
			if err != nil {
				return nil, err
			}

			if len(missedData) > 0 {
				for _, data := range missedData {
					itemMap[data.ID] = data
				}
				err = d.cache.MultiSet(ctx, missedData, cache.UsersExpireTime)
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
func (d *usersDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.Users, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.Users{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (d *usersDao) GetByEmail(ctx context.Context, params string) (*model.Users, error) {
	user := &model.Users{}
	err := d.db.WithContext(ctx).Where("email = ?", params).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByEmailTx GetByEmail query by email
func (d *usersDao) GetByEmailTx(ctx context.Context, tx *gorm.DB, params *model.Users) (*model.Users, error) {
	user := &model.Users{}
	err := tx.WithContext(ctx).Where("email = ?", params.Email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *usersDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Users) (uint64, error) {
	if tx == nil {
		return 0, errors.New("transaction is nil")
	}
	if table == nil {
		return 0, errors.New("user model is nil")
	}
	if ctx == nil {
		return 0, errors.New("context is nil")
	}
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *usersDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.Users{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *usersDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Users) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
