package consume

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/socket"
)

// 撤销聊天消息
func (h *IMHandler) onConsumeTalkRevoke(ctx context.Context, body []byte) {
	var in types.ConsumeTalkRevoke
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkRevoke Unmarshal err: %s", err.Error())
		return
	}

	var record model.TalkRecords
	if err := h.db.First(&record, "msg_id = ?", in.MsgId).Error; err != nil {
		return
	}

	var clientIds []int64
	if record.TalkType == constant.ChatPrivateMode {
		for _, uid := range [2]int{record.UserID, record.ReceiverID} {
			ids := h.messageCache.GetUidFromClientIds(ctx, h.config.App.Sid, socket.Session.Chat.Name(), strconv.Itoa(uid))
			clientIds = append(clientIds, ids...)
		}
	} else if record.TalkType == constant.ChatGroupMode {
		clientIds = h.chatRoom.All(ctx, &types.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: constant.RoomImGroup,
			Number:   strconv.Itoa(record.ReceiverID),
			Sid:      h.config.App.Sid,
		})
	}

	if len(clientIds) == 0 {
		return
	}

	var user model.Users
	if err := h.db.WithContext(ctx).Select("id,nickname").First(&user, record.UserID).Error; err != nil {
		return
	}

	c := socket.NewSenderContent()
	c.SetAck(true)
	c.SetReceive(clientIds...)
	c.SetMessage(constant.PushEventImMessageRevoke, map[string]any{
		"talk_type":   record.TalkType,
		"sender_id":   record.UserID,
		"receiver_id": record.ReceiverID,
		"msg_id":      record.MsgID,
		"text":        fmt.Sprintf("%s: 撤回了一条消息", user.Username),
	})

	socket.Session.Chat.Write(c)
}
