package dao

import (
	"context"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
)

const (
	GroupMemberQuitStatusYes = 1
	GroupMemberQuitStatusNo  = 0

	GroupMemberMuteStatusYes = 1
	GroupMemberMuteStatusNo  = 0
)

type GroupMemberDao interface {
	IsMaster(ctx context.Context, gid, uid int) bool
	IsLeader(ctx context.Context, gid, uid int) bool
	IsMember(ctx context.Context, gid, uid int, cache bool) bool
	FindByUserId(ctx context.Context, gid, uid int) (*model.GroupMember, error)
	GetMemberIds(ctx context.Context, groupId int) []int
	GetUserGroupIds(ctx context.Context, uid int) []int
	CountMemberTotal(ctx context.Context, gid int) int64
	GetMemberRemark(ctx context.Context, groupId int, userId int) string
	GetMembers(ctx context.Context, groupId int) []*types.MemberItem
	CountGroupMemberNum(ids []int) ([]*model.CountGroupMember, error)
	CheckUserGroup(ids []int, userId int) ([]int, error)
}

type groupMemberDao struct {
	db    *gorm.DB
	cache cache.GroupMemberCache // if nil, the cache is not used.
	sfg   *singleflight.Group    // if cache is nil, the sfg is not used.
}

func NewGroupMemberDao(db *gorm.DB, xCache cache.GroupMemberCache) GroupMemberDao {
	if xCache == nil {
		return &groupMemberDao{db: db}
	}
	return &groupMemberDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (g groupMemberDao) IsMaster(ctx context.Context, gid, uid int) bool {
	exist, err := g.queryExist(ctx, "group_id = ? and user_id = ? and leader = 2 and is_quit = ?", gid, uid, GroupMemberQuitStatusNo)
	return err == nil && exist
}

func (g groupMemberDao) IsLeader(ctx context.Context, gid, uid int) bool {
	exist, err := g.queryExist(ctx, "group_id = ? and user_id = ? and leader in (1,2) and is_quit = ?", gid, uid, GroupMemberQuitStatusNo)
	return err == nil && exist
}

func (g groupMemberDao) IsMember(ctx context.Context, gid, uid int, cache bool) bool {
	if cache && g.cache.IsGroupRelation(ctx, uid, gid) == nil {
		return true
	}

	exist, err := g.queryExist(ctx, "group_id = ? and user_id = ? and is_quit = ?", gid, uid, GroupMemberQuitStatusNo)
	if err != nil {
		return false
	}

	if exist {
		g.cache.SetGroupRelation(ctx, uid, gid)
	}

	return exist
}

func (g groupMemberDao) FindByUserId(ctx context.Context, gid, uid int) (*model.GroupMember, error) {
	member := &model.GroupMember{}
	err := g.db.Model(ctx).Where("group_id = ? and user_id = ?", gid, uid).First(member).Error
	return member, err
}

func (g groupMemberDao) GetMemberIds(ctx context.Context, groupId int) []int {
	var ids []int
	_ = g.db.Model(ctx).Select("user_id").Where("group_id = ? and is_quit = ?", groupId, GroupMemberQuitStatusNo).Scan(&ids)

	return ids
}

func (g groupMemberDao) GetUserGroupIds(ctx context.Context, uid int) []int {
	var ids []int
	_ = g.db.Model(ctx).Where("user_id = ? and is_quit = ?", uid, GroupMemberQuitStatusNo).Pluck("group_id", &ids)

	return ids
}

func (g groupMemberDao) CountMemberTotal(ctx context.Context, gid int) int64 {
	count, _ := g.queryCount(ctx, "group_id = ? and is_quit = ?", gid, GroupMemberQuitStatusNo)
	return count
}

func (g groupMemberDao) GetMemberRemark(ctx context.Context, groupId int, userId int) string {
	var remarks string
	g.db.Model(ctx).Select("user_card").Where("group_id = ? and user_id = ?", groupId, userId).Scan(&remarks)

	return remarks
}

func (g groupMemberDao) GetMembers(ctx context.Context, groupId int) []*types.MemberItem {
	fields := []string{
		"group_member.id",
		"group_member.leader",
		"group_member.user_card",
		"group_member.user_id",
		"group_member.is_mute",
		"users.avatar",
		"users.nickname",
		"users.gender",
		"users.motto",
	}

	tx := g.db.WithContext(ctx).Table("group_member")
	tx.Joins("left join users on users.id = group_member.user_id")
	tx.Where("group_member.group_id = ? and group_member.is_quit = ?", groupId, GroupMemberQuitStatusNo)
	tx.Order("group_member.leader desc")

	var items []*types.MemberItem
	tx.Unscoped().Select(fields).Scan(&items)

	return items
}

func (g groupMemberDao) CountGroupMemberNum(ids []int) ([]*model.CountGroupMember, error) {
	var items []*model.CountGroupMember
	err := g.db.Model(context.TODO()).Select("group_id,count(*) as count").Where("group_id in ? and is_quit = ?", ids, GroupMemberQuitStatusNo).Group("group_id").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (g groupMemberDao) CheckUserGroup(ids []int, userId int) ([]int, error) {
	items := make([]int, 0)

	err := g.db.Model(context.TODO()).Select("group_id").Where("group_id in ? and user_id = ? and is_quit = ?", ids, userId, GroupMemberQuitStatusNo).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (g groupMemberDao) queryExist(ctx context.Context, where string, args ...any) (bool, error) {

	var count int64
	err := g.db.Model(ctx).Select("1").Where(where, args...).Limit(1).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

// QueryCount 根据条件统计数据总数
func (g groupMemberDao) queryCount(ctx context.Context, where string, args ...any) (int64, error) {

	var count int64
	err := g.db.Model(ctx).Where(where, args...).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
