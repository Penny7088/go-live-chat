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

// 消息已读事件
func (h *IMHandler) onConsumeTalkRead(ctx context.Context, body []byte) {
	var in types.ConsumeTalkRead
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeContactApply Unmarshal err: %s", err.Error())
		return
	}

	clientIds := h.messageCache.GetUidFromClientIds(ctx, h.config.App.Sid, socket.Session.Chat.Name(), strconv.Itoa(in.ReceiverId))
	if len(clientIds) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetAck(true)
	c.SetReceive(clientIds...)
	c.SetMessage(constant.PushEventImMessageRead, map[string]any{
		"sender_id":   in.SenderId,
		"receiver_id": in.ReceiverId,
		"msg_ids":     in.MsgIds,
	})

	socket.Session.Chat.Write(c)
}
