package handler

import "lingua_exchange/internal/cache"

var _ MessageHandler = (*messageHandler)(nil)

type MessageHandler interface {
}

type messageHandler struct {
	messageCache cache.MessageCache
}
