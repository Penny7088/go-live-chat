package handler

import "context"

var _ TalkRecordsHandler = (*talkRecordsHandler)(nil)

type TalkRecordsHandler interface {
}

type talkRecordsHandler struct {
}

func (t talkRecordsHandler) Publish(ctx context.Context) {
	// TODO implement me
	panic("implement me")
}
