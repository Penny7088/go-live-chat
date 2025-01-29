package imService

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
)

// 消息权限控制类
var _ IPermissionService = (*permissionService)(nil)

type IPermissionService interface {
	IsAuth(ctx context.Context, opt *types.AuthOption) error
}

type permissionService struct {
	groupDao       dao.GroupDao
	groupMemberDao dao.GroupMemberDao
}

func NewPermissionService() IPermissionService {
	return &permissionService{
		groupDao:       dao.NewGroupDao(model.GetDB(), cache.NewGroupCache(model.GetCacheType())),
		groupMemberDao: dao.NewGroupMemberDao(model.GetDB(), cache.NewGroupMemberCache(model.GetCacheType())),
	}
}

func (p permissionService) IsAuth(ctx context.Context, opt *types.AuthOption) error {
	// / todo 私聊目前不做任何验证 后期考虑是否需要
	// if opt.TalkType == constant.ChatPrivateMode{
	//
	// }

	group, err := p.groupDao.GetByID(ctx, opt.ReceiverId)
	if err != nil {
		return err
	}
	if group.IsDismiss == 1 {
		return errors.New("此群聊已解散")
	}

	memberInfo, err := p.groupMemberDao.FindByUserId(ctx, int(opt.ReceiverId), opt.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("暂无权限发送消息！")
		}

		return errors.New("系统繁忙，请稍后再试！！！")
	}
	if memberInfo.IsQuit == constant.GroupMemberQuitStatusYes {
		return errors.New("暂无权限发送消息！")
	}

	if memberInfo.IsMute == constant.GroupMemberMuteStatusYes {
		return errors.New("已被群主或管理员禁言！")
	}
	return nil
}
