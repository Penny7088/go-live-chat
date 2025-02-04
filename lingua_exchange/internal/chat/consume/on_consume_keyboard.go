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

// 键盘输入事件消息
func (h *IMHandler) onConsumeTalkKeyboard(ctx context.Context, body []byte) {

	var in types.ConsumeTalkKeyboard
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkKeyboard Unmarshal err: %s", err.Error())
		return
	}

	ids := h.messageCache.GetUidFromClientIds(ctx, h.config.App.Sid, socket.Session.Chat.Name(), strconv.Itoa(in.ReceiverID))
	if len(ids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(ids...)
	c.SetMessage(constant.PushEventImMessageKeyboard, map[string]any{
		"sender_id":   in.SenderID,
		"receiver_id": in.ReceiverID,
	})

	socket.Session.Chat.Write(c)
}
