package consume

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/socket"
)

// 聊天消息事件
func (h *IMHandler) onConsumeTalk(ctx context.Context, body []byte) {

	var in types.ConsumeTalk
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalk Unmarshal err: %s", err.Error())
		return
	}

	var clientIds []int64
	if in.TalkType == constant.ChatPrivateMode {
		for _, val := range [2]int64{in.SenderID, in.ReceiverID} {
			ids := h.messageCache.GetUidFromClientIds(ctx, h.config.App.Sid, socket.Session.Chat.Name(), strconv.FormatInt(val, 10))

			clientIds = append(clientIds, ids...)
		}
	} else if in.TalkType == constant.ChatGroupMode {
		ids := h.chatRoom.All(ctx, &types.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: constant.RoomImGroup,
			Number:   strconv.Itoa(int(in.ReceiverID)),
			Sid:      h.config.App.Sid,
		})

		clientIds = append(clientIds, ids...)
	}

	if len(clientIds) == 0 {
		return
	}

	data, err := h.talkRecordsDao.FindTalkRecord(ctx, in.MsgId)
	if err != nil {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetAck(true)
	c.SetMessage(constant.PushEventImMessage, map[string]any{
		"sender_id":   in.SenderID,
		"receiver_id": in.ReceiverID,
		"talk_type":   in.TalkType,
		"data":        data,
	})

	socket.Session.Chat.Write(c)
}
