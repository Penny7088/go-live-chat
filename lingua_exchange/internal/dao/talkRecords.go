package dao

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/jsonutil"
	"lingua_exchange/pkg/sliceutil"
	"lingua_exchange/pkg/timeutil"

	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"

	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/model"
)

var _ TalkRecordsDao = (*talkRecordsDao)(nil)

// TalkRecordsDao defining the dao interface
type TalkRecordsDao interface {
	Create(ctx context.Context, table *model.TalkRecords) error
	DeleteByID(ctx context.Context, id string) error
	UpdateByID(ctx context.Context, table *model.TalkRecords) error
	GetByID(ctx context.Context, id string) (*model.TalkRecords, error)

	DeleteByIDs(ctx context.Context, ids []string) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.TalkRecords, error)
	GetByIDs(ctx context.Context, ids []string) (map[string]*model.TalkRecords, error)
	GetByLastID(ctx context.Context, lastID string, limit int, sort string) ([]*model.TalkRecords, error)

	DeleteByTx(ctx context.Context, tx *gorm.DB, id string) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.TalkRecords) error
	FindTalkRecord(ctx context.Context, id string) (*types.TalkRecordItem, error)
	FindAllTalkRecords(ctx context.Context, opt *types.FindAllTalkRecordsOpt) ([]*types.TalkRecordItem, error)
}

type talkRecordsDao struct {
	db                   *gorm.DB
	cache                cache.TalkRecordsCache // if nil, the cache is not used.
	sfg                  *singleflight.Group    // if cache is nil, the sfg is not used.
	talkRecordsVoteDao   TalkRecordsVoteDao
	talkRecordsDeleteDao TalkRecordsDeleteDao
}

// NewTalkRecordsDao creating the dao interface
func NewTalkRecordsDao(db *gorm.DB, xCache cache.TalkRecordsCache) TalkRecordsDao {
	if xCache == nil {
		return &talkRecordsDao{db: db}
	}
	return &talkRecordsDao{
		db:                   db,
		cache:                xCache,
		sfg:                  new(singleflight.Group),
		talkRecordsVoteDao:   NewTalkRecordsVoteDao(db, cache.NewTalkRecordsVoteCache(model.GetCacheType())),
		talkRecordsDeleteDao: NewTalkRecordsDeleteDao(db, cache.NewTalkRecordsDeleteCache(model.GetCacheType())),
	}
}

func (d *talkRecordsDao) FindAllTalkRecords(ctx context.Context, opt *types.FindAllTalkRecordsOpt) ([]*types.TalkRecordItem, error) {
	var (
		items  = make([]*types.QueryTalkRecord, 0, opt.Limit)
		cursor = opt.Cursor
	)

	for {
		// 这里查询数据放弃了关联查询，所以这里需要查询多次，防止查询中存在用户已删除的数据需要过滤掉
		list, err := d.findAllRecords(ctx, &types.FindAllTalkRecordsOpt{
			TalkType:   opt.TalkType,
			UserId:     opt.UserId,
			ReceiverId: opt.ReceiverId,
			MsgType:    opt.MsgType,
			Cursor:     cursor,
			Limit:      opt.Limit + 10, // 多查几条数据
		})

		if err != nil {
			return nil, err
		}

		if len(list) == 0 {
			break
		}

		tmpMsgIds := make([]string, 0, len(list))
		for _, v := range list {
			tmpMsgIds = append(tmpMsgIds, v.MsgId)
		}

		msgIds, err := d.talkRecordsDeleteDao.FindAllMsgIds(ctx, tmpMsgIds, opt.UserId)
		if err != nil {
			return nil, err
		}

		hashIds := make(map[string]struct{}, len(msgIds))
		for _, msgId := range msgIds {
			hashIds[msgId] = struct{}{}
		}

		for _, v := range list {
			if _, ok := hashIds[v.MsgId]; ok {
				continue
			}

			items = append(items, v)
		}

		if len(items) >= opt.Limit || len(list) < opt.Limit {
			break
		}

		// 设置游标继续往下执行
		cursor = int(list[len(list)-1].Sequence)
	}

	if len(items) > opt.Limit {
		items = items[:opt.Limit]
	}

	return d.handleTalkRecords(ctx, items)

}

func (d *talkRecordsDao) findAllRecords(ctx context.Context, opt *types.FindAllTalkRecordsOpt) ([]*types.QueryTalkRecord, error) {
	query := d.db.WithContext(ctx).Table("talk_records")
	query.Select([]string{
		"talk_records.sequence",
		"talk_records.talk_type",
		"talk_records.msg_type",
		"talk_records.msg_id",
		"talk_records.user_id",
		"talk_records.receiver_id",
		"talk_records.is_revoke",
		"talk_records.extra",
		"talk_records.created_at",
	})

	if opt.Cursor > 0 {
		query.Where("talk_records.sequence < ?", opt.Cursor)
	}

	if opt.TalkType == constant.ChatPrivateMode {
		subQuery := d.db.Where("talk_records.user_id = ? and talk_records.receiver_id = ?", opt.UserId, opt.ReceiverId)
		subQuery.Or("talk_records.user_id = ? and talk_records.receiver_id = ?", opt.ReceiverId, opt.UserId)

		query.Where(subQuery)
	} else {
		query.Where("talk_records.receiver_id = ?", opt.ReceiverId)
	}

	if opt.MsgType != nil && len(opt.MsgType) > 0 {
		query.Where("talk_records.msg_type in ?", opt.MsgType)
	}

	query.Where("talk_records.talk_type = ?", opt.TalkType)
	query.Order("talk_records.sequence desc").Limit(opt.Limit)

	var items []*types.QueryTalkRecord
	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (d *talkRecordsDao) deleteCache(ctx context.Context, id string) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *talkRecordsDao) Create(ctx context.Context, table *model.TalkRecords) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *talkRecordsDao) DeleteByID(ctx context.Context, id string) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.TalkRecords{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *talkRecordsDao) UpdateByID(ctx context.Context, table *model.TalkRecords) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.MsgID)

	return err
}

func (d *talkRecordsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.TalkRecords) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.MsgID != "" {
		update["msg_id"] = table.MsgID
	}
	if table.Sequence != 0 {
		update["sequence"] = table.Sequence
	}
	if table.TalkType != 0 {
		update["talk_type"] = table.TalkType
	}
	if table.MsgType != 0 {
		update["msg_type"] = table.MsgType
	}
	if table.UserID != 0 {
		update["user_id"] = table.UserID
	}
	if table.ReceiverID != 0 {
		update["receiver_id"] = table.ReceiverID
	}
	if table.IsRevoke != 0 {
		update["is_revoke"] = table.IsRevoke
	}
	if table.IsMark != 0 {
		update["is_mark"] = table.IsMark
	}
	if table.QuoteID != "" {
		update["quote_id"] = table.QuoteID
	}
	if table.Extra != "" {
		update["extra"] = table.Extra
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *talkRecordsDao) GetByID(ctx context.Context, id string) (*model.TalkRecords, error) {
	// no cache
	if d.cache == nil {
		record := &model.TalkRecords{}
		err := d.db.WithContext(ctx).Where("msg_id = ?", id).First(record).Error
		return record, err
	}

	// get from cache or database
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	if errors.Is(err, model.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(id, func() (interface{}, error) { // nolint
			table := &model.TalkRecords{}
			err = d.db.WithContext(ctx).Where("msg_id = ?", id).First(table).Error
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
			err = d.cache.Set(ctx, id, table, cache.TalkRecordsExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, msg_id=%s", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.TalkRecords)
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
func (d *talkRecordsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.TalkRecords, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.TalkRecords{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.TalkRecords{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs delete records by batch id
func (d *talkRecordsDao) DeleteByIDs(ctx context.Context, ids []string) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.TalkRecords{}).Error
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
func (d *talkRecordsDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.TalkRecords, error) {
	queryStr, args, err := c.ConvertToGorm()
	if err != nil {
		return nil, err
	}

	table := &model.TalkRecords{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs get records by batch id
func (d *talkRecordsDao) GetByIDs(ctx context.Context, ids []string) (map[string]*model.TalkRecords, error) {
	// no cache
	if d.cache == nil {
		var records []*model.TalkRecords
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[string]*model.TalkRecords)
		for _, record := range records {
			itemMap[record.MsgID] = record
		}
		return itemMap, nil
	}

	// get form cache or database
	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	var missedIDs []string
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
		var realMissedIDs []string
		for _, id := range missedIDs {
			_, err = d.cache.Get(ctx, id)
			if errors.Is(err, cacheBase.ErrPlaceholder) {
				continue
			}
			realMissedIDs = append(realMissedIDs, id)
		}

		if len(realMissedIDs) > 0 {
			var missedData []*model.TalkRecords
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
			if err != nil {
				return nil, err
			}

			if len(missedData) > 0 {
				for _, data := range missedData {
					itemMap[data.MsgID] = data
				}
				err = d.cache.MultiSet(ctx, missedData, cache.TalkRecordsExpireTime)
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
func (d *talkRecordsDao) GetByLastID(ctx context.Context, lastID string, limit int, sort string) ([]*model.TalkRecords, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.TalkRecords{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *talkRecordsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.TalkRecords) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *talkRecordsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id string) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.TalkRecords{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *talkRecordsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.TalkRecords) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.MsgID)

	return err
}

func (d *talkRecordsDao) FindTalkRecord(ctx context.Context, id string) (*types.TalkRecordItem, error) {
	var (
		err    error
		item   *types.QueryTalkRecord
		fields = []string{
			"talk_records.msg_id",
			"talk_records.sequence",
			"talk_records.talk_type",
			"talk_records.msg_type",
			"talk_records.user_id",
			"talk_records.receiver_id",
			"talk_records.is_revoke",
			"talk_records.extra",
			"talk_records.created_at",
		}
	)

	query := d.db.Table("talk_records")
	query.Where("talk_records.msg_id = ?", id)
	if err = query.Select(fields).Take(&item).Error; err != nil {
		return nil, err
	}

	list, err := d.handleTalkRecords(ctx, []*types.QueryTalkRecord{item})
	if err != nil {
		return nil, err
	}

	return list[0], nil
}

// HandleTalkRecords 处理消息
func (s *talkRecordsDao) handleTalkRecords(ctx context.Context, items []*types.QueryTalkRecord) ([]*types.TalkRecordItem, error) {
	if len(items) == 0 {
		return make([]*types.TalkRecordItem, 0), nil
	}

	var (
		votes     []string
		voteItems []*model.TalkRecordsVote
	)

	uids := make([]int, 0, len(items))
	for _, item := range items {
		uids = append(uids, item.UserId)

		switch item.MsgType {
		case constant.ChatMsgTypeVote:
			votes = append(votes, item.MsgId)
		}
	}

	var usersItems []*model.Users
	err := s.db.Model(&model.Users{}).Select("id,username,profile_picture").Where("id in ?", sliceutil.Unique(uids)).Scan(&usersItems).Error
	if err != nil {
		return nil, err
	}

	hashUser := make(map[uint64]*model.Users)
	for _, user := range usersItems {
		hashUser[user.ID] = user
	}

	hashVotes := make(map[string]*model.TalkRecordsVote)
	if len(votes) > 0 {
		s.db.Model(&model.TalkRecordsVote{}).Where("msg_id in ?", votes).Scan(&voteItems)
		for i := range voteItems {
			hashVotes[voteItems[i].MsgID] = voteItems[i]
		}
	}

	newItems := make([]*types.TalkRecordItem, 0, len(items))
	for _, item := range items {
		data := &types.TalkRecordItem{
			MsgId:      item.MsgId,
			Sequence:   int(item.Sequence),
			TalkType:   item.TalkType,
			MsgType:    item.MsgType,
			UserId:     item.UserId,
			ReceiverId: item.ReceiverId,
			Nickname:   item.Nickname,
			Avatar:     item.Avatar,
			IsRevoke:   item.IsRevoke,
			IsMark:     item.IsMark,
			CreatedAt:  timeutil.FormatDatetime(item.CreatedAt),
			Extra:      make(map[string]any),
		}

		if user, ok := hashUser[uint64(item.UserId)]; ok {
			data.Nickname = user.Username
			data.Avatar = user.ProfilePicture
		}

		_ = jsonutil.Decode(item.Extra, &data.Extra)

		switch item.MsgType {
		case constant.ChatMsgTypeVote:
			if value, ok := hashVotes[item.MsgId]; ok {
				options := make(map[string]any)
				opts := make([]any, 0)

				if err := jsonutil.Decode(value.AnswerOption, &options); err == nil {
					arr := make([]string, 0, len(options))
					for k := range options {
						arr = append(arr, k)
					}

					sort.Strings(arr)

					for _, v := range arr {
						opts = append(opts, map[string]any{
							"key":   v,
							"value": options[v],
						})
					}
				}

				users := make([]int, 0)
				if uids, err := s.talkRecordsVoteDao.GetVoteAnswerUser(ctx, int(value.ID)); err == nil {
					users = uids
				}

				var statistics any

				if res, err := s.talkRecordsVoteDao.GetVoteStatistics(ctx, int(value.ID)); err != nil {
					statistics = map[string]any{
						"count":   0,
						"options": map[string]int{},
					}
				} else {
					statistics = res
				}

				data.Extra = map[string]any{
					"detail": map[string]any{
						"id":            value.ID,
						"msg_id":        value.MsgID,
						"title":         value.Title,
						"answer_mode":   value.AnswerMode,
						"status":        value.Status,
						"answer_option": opts,
						"answer_num":    value.AnswerNum,
						"answered_num":  value.AnsweredNum,
					},
					"statistics": statistics,
					"vote_users": users, // 已投票成员
				}
			}
		}

		newItems = append(newItems, data)
	}

	return newItems, nil
}
