package consume

import (
	"context"
	"encoding/json"

	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/types"
)

// 用户上线或下线消息
func (h *IMHandler) onConsumeContactStatus(ctx context.Context, body []byte) {

	var in types.ConsumeContactStatus
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeContactStatus Unmarshal err: %s", err.Error())
		return
	}
}

// 好友申请消息
func (h *IMHandler) onConsumeContactApply(ctx context.Context, body []byte) {

}
