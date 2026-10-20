package mocks

import "testing"

type RabbitClientMock struct {
	receiver  *RabbitMqReceiver
	publisher *RabbitMqSender
	t         *testing.T
}

func NewRabbitClientMock(t *testing.T) *RabbitClientMock {
	return &RabbitClientMock{t: t}
}

func (r *RabbitClientMock) GetReceiver(expectedErrorConsume error) (*RabbitMqReceiver, error) {
	if r.receiver == nil {
		r.receiver = NewRabbitMqReceiver(r.t, expectedErrorConsume)
	}
	return r.receiver, nil
}

func (r *RabbitClientMock) GetSender(expectedErrorPush error) (*RabbitMqSender, error) {
	if r.receiver == nil {
		r.publisher = NewRabbitMqSender(r.t, expectedErrorPush)
	}
	return r.publisher, nil
}
