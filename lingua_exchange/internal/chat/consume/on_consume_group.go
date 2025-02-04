package consume

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/socket"
)

// 加入群房间
func (h *IMHandler) onConsumeGroupJoin(ctx context.Context, body []byte) {

	var in types.ConsumeGroupJoin
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeGroupJoin Unmarshal err: %s", err.Error())
		return
	}

	sid := h.config.App.Sid
	for _, uid := range in.Uids {
		ids := h.messageCache.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), strconv.Itoa(uid))

		for _, cid := range ids {
			opt := &types.RoomOption{
				Channel:  socket.Session.Chat.Name(),
				RoomType: constant.RoomImGroup,
				Number:   strconv.Itoa(in.Gid),
				Sid:      h.config.App.Sid,
				Cid:      cid,
			}

			if in.Type == 2 {
				_ = h.chatRoom.Del(ctx, opt)
			} else {
				_ = h.chatRoom.Add(ctx, opt)
			}
		}
	}
}

// 入群申请通知
func (h *IMHandler) onConsumeGroupApply(ctx context.Context, body []byte) {

	var in types.ConsumeGroupApply
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeGroupApply Unmarshal err: %s", err.Error())
		return
	}

	var groupMember model.GroupMember
	if err := h.db.First(&groupMember, "group_id = ? and leader = ?", in.GroupId, 2).Error; err != nil {
		return
	}

	var groupDetail model.Group
	if err := h.db.First(&groupDetail, in.GroupId).Error; err != nil {
		return
	}

	var user model.Users
	if err := h.db.First(&user, in.UserId).Error; err != nil {
		return
	}

	data := make(map[string]any)
	data["group_name"] = groupDetail.Name
	data["username"] = user.Username

	clientIds := h.messageCache.GetUidFromClientIds(ctx, h.config.App.Sid, socket.Session.Chat.Name(), strconv.Itoa(groupMember.UserID))

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetMessage(constant.PushEventGroupApply, data)

	socket.Session.Chat.Write(c)
}
