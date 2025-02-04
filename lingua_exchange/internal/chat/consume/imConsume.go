package consume

type IMConsumer struct {
}

func NewIMConsumer() *IMConsumer {
	return &IMConsumer{}
}

func (I IMConsumer) Call(event string, data []byte) {
	// TODO implement me
	panic("implement me")
}
