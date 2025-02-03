package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/imService"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/jsonutil"
	"lingua_exchange/pkg/jwt"
	"lingua_exchange/pkg/timeutil"
)

var _ GroupNoticeHandler = (*groupNoticeHandler)(nil)

type GroupNoticeHandler interface {
	List(ctx *gin.Context)
	CreateAndUpdate(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type groupNoticeHandler struct {
	iDao           dao.GroupNoticeDao
	messageService imService.IMessageService
	groupMemberDao dao.GroupMemberDao
}

func NewGroupNoticeHandler() GroupNoticeHandler {
	return &groupNoticeHandler{
		iDao:           dao.NewGroupNoticeDao(model.GetDB(), cache.NewGroupNoticeCache(model.GetCacheType())),
		messageService: imService.NewMessageService(),
		groupMemberDao: dao.NewGroupMemberDao(model.GetDB(), cache.NewGroupMemberCache(model.GetCacheType())),
	}
}

func (g groupNoticeHandler) List(c *gin.Context) {
	params := &types.GroupNoticeListRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	uid, err2 := jwt.HeaderObtainUID(c)
	if err2 != nil {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	if !g.groupMemberDao.IsMember(ctx, params.GroupID, uid, true) {
		logger.Warn("无获取数据权限 ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Forbidden)
		return
	}
	all, err := g.iDao.GetListAll(ctx, params.GroupID)
	if err != nil {
		logger.Warn("获取所有公告失败 ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrAllGroupNotice)
		return
	}

	items := make([]*types.NoticeItemReply, 0)
	for i := 0; i < len(all); i++ {
		items = append(items, &types.NoticeItemReply{
			ID:           int32(all[i].Id),
			Title:        all[i].Title,
			Content:      all[i].Content,
			IsTop:        int32(all[i].IsTop),
			IsConfirm:    int32(all[i].IsConfirm),
			ConfirmUsers: all[i].ConfirmUsers,
			Avatar:       all[i].Avatar,
			CreatorID:    int32(all[i].CreatorId),
			CreatedAt:    timeutil.FormatDatetime(all[i].CreatedAt),
			UpdatedAt:    timeutil.FormatDatetime(all[i].UpdatedAt),
		})
	}

	response.Success(c, gin.H{
		"list": items,
	})
}

func (g groupNoticeHandler) CreateAndUpdate(c *gin.Context) {
	params := &types.GroupNoticeEditRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	uid, err2 := jwt.HeaderObtainUID(c)
	if err2 != nil {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	if !g.groupMemberDao.IsLeader(ctx, params.GroupID, uid) {
		logger.Warn("无权限操作 ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}

	err, msg := g.updateOrCreate(params, ctx, uid)
	if err != nil {
		logger.Warn("创建或者更新失败 ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrNoticeCreateOrUpdateFailed)
		return
	}

	_ = g.messageService.SendSysOther(&ctx, &model.TalkRecords{
		TalkType:   constant.TalkRecordTalkTypeGroup,
		MsgType:    constant.ChatMsgTypeGroupNotice,
		UserID:     uid,
		ReceiverID: params.GroupID,
		Extra: jsonutil.Encode(types.TalkRecordExtraGroupNotice{
			OwnerId:   uid,
			OwnerName: "owner",
			Title:     params.Title,
			Content:   params.Content,
		}),
	})

	response.Success(c, gin.H{
		"message": msg,
	})

}

func (g groupNoticeHandler) updateOrCreate(params *types.GroupNoticeEditRequest, ctx context.Context, uid int) (error, string) {
	if params.NoticeID == 0 {
		err := g.iDao.Create(ctx, &model.GroupNotice{
			CreatorID:    uint(uid),
			GroupID:      uint(params.GroupID),
			Title:        params.Title,
			Content:      params.Content,
			IsTop:        uint(params.IsTop),
			IsConfirm:    uint(params.IsConfirm),
			ConfirmUsers: "{}",
		})
		if err != nil {
			return err, ""
		}
		return nil, "添加群公告成功"
	} else {
		err := g.iDao.UpdateByID(ctx, &model.GroupNotice{
			GroupID:   uint(params.GroupID),
			ID:        uint64(params.NoticeID),
			Title:     params.Title,
			Content:   params.Content,
			IsTop:     uint(params.IsTop),
			IsConfirm: uint(params.IsConfirm),
		})
		if err != nil {
			return err, ""
		}
		return nil, "更新群公告成功"
	}
}

func (g groupNoticeHandler) Delete(c *gin.Context) {
	params := &types.GroupNoticeDeleteRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := g.iDao.UpdateByID(ctx, &model.GroupNotice{
		IsDelete:  1,
		DeletedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uint64(params.NoticeID),
		GroupID:   uint(params.GroupID),
	})

	if err != nil {
		logger.Warn("群公告删除成功 ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNoticeDeleteFailed)
		return
	}

	response.Success(c, "ok")
}
