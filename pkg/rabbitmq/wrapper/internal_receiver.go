package wrapper

import (
	"context"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

// 1. Định nghĩa Interface để Mocking trong file test
//
//go:generate mockery --name InternalConsumer --filename consumer_mock.go --structname MockInternalConsumer
type InternalConsumer interface {
	Receive(ctx context.Context) (AmqpDelivery, error)
	Close(ctx context.Context) error
}

// realConsumerWrapper wrap *rmq.Consumer and implements *rmq.Consumer
type realConsumerWrapper struct {
	realConsumer *rmq.Consumer
}

func NewRabbitReceiverInstance(realConsumer *rmq.Consumer) InternalConsumer {
	return &realConsumerWrapper{realConsumer: realConsumer}
}

// Receive message rmpq.Consumer
//
// Parameters:
//   - ctx: The context for the publisher.
//
// Returns:
//   - error
func (w *realConsumerWrapper) Receive(ctx context.Context) (AmqpDelivery, error) {
	d, err := w.realConsumer.Receive(ctx)
	if err != nil {
		return nil, err
	}
	// wrap delivery
	return NewDeliveryWrapper(d), nil
}

// Close rmpq.Consumer
//
// Parameters:
//   - ctx: The context for the publisher.
//
// Returns:
//   - error
func (w *realConsumerWrapper) Close(ctx context.Context) error {
	return w.realConsumer.Close(ctx)
}
