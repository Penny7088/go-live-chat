package dao

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"lingua_exchange/internal/model"
	"lingua_exchange/pkg/timeutil"
)

var _ TalkSessionDao = (*talkSessionDao)(nil)

type TalkSessionDao interface {
	IsDisturb(ctx context.Context, uid int, receiverId int, talkType int) bool
	FindBySessionId(ctx context.Context, uid int, receiverId int, talkType int) int
	BatchAddList(ctx context.Context, uid int, values map[string]int)
	List(ctx context.Context, uid int) ([]*model.SearchTalkSession, error)
}

type talkSessionDao struct {
	db *gorm.DB
}

func (t talkSessionDao) List(ctx context.Context, uid int) ([]*model.SearchTalkSession, error) {
	fields := []string{
		"list.id", "list.talk_type", "list.receiver_id", "list.updated_at",
		"list.is_disturb", "list.is_top", "list.is_robot",
		"`users`.profile_picture as user_avatar", "`users`.username",
		"`group`.name as group_name", "`group`.avatar as group_avatar",
	}
	query := t.db.WithContext(ctx).Table("talk_sessions list")
	query.Joins("left join `users` ON list.receiver_id = `users`.id AND list.talk_type = 1")
	query.Joins("left join `group` ON list.receiver_id = `group`.id AND list.talk_type = 2")
	query.Where("list.user_id = ? and list.is_delete = 0", uid)
	query.Order("list.updated_at desc")

	var items []*model.SearchTalkSession
	if err := query.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func NewTalkSessionDao(db *gorm.DB) TalkSessionDao {

	return &talkSessionDao{
		db: db,
	}
}

func (t talkSessionDao) BatchAddList(ctx context.Context, uid int, values map[string]int) {
	ctime := timeutil.DateTime()
	data := make([]string, 0)
	for k, v := range values {
		if v == 0 {
			continue
		}

		value := strings.Split(k, "_")
		if len(value) != 2 {
			continue
		}

		data = append(data, fmt.Sprintf("(%s, %d, %s, '%s', '%s')", value[0], uid, value[1], ctime, ctime))
	}
	if len(data) == 0 {
		return
	}
	sprintf := fmt.Sprintf("INSERT INTO talk_session ( `talk_type`, `user_id`, `receiver_id`, created_at, updated_at ) VALUES %s ON DUPLICATE KEY UPDATE is_delete = 0, updated_at = '%s'", strings.Join(data, ","), ctime)
	t.db.WithContext(ctx).Exec(sprintf)
}

func (t talkSessionDao) IsDisturb(ctx context.Context, uid int, receiverId int, talkType int) bool {
	return false
}

func (t talkSessionDao) FindBySessionId(ctx context.Context, uid int, receiverId int, talkType int) int {
	// TODO implement me
	panic("implement me")
}
