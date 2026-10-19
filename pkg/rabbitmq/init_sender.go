package rabbitmq

import (
	"context"
	"fmt"

	internal "github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq/wrapper"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

//go:generate mockery --name RabbitMqSender --filename sender.go
type RabbitMqSender interface {
	Publish(ctx context.Context, message []byte) error
	Close() error
}

type rabbitMqSender struct {
	ctx          context.Context
	publisherRmq internal.InternalPublisher // Thay đổi từ *rmq.Publisher thành interface
}

// NewRabbitSenderInstance creates a new rabbitMqSender instance with the provided publisher.
//
// Parameters:
//   - ctx: The context for the publisher.
//   - publisher: The publisher to use for publishing messages, which stands for "rabbitmq.Publisher".
//
// Returns:
//   - An RabbitMqSender
//
// Warning:
//   - Never pass publisher parameter is nil.
func NewRabbitSenderInstance(ctx context.Context, publisher internal.InternalPublisher) RabbitMqSender {
	return &rabbitMqSender{ctx: ctx, publisherRmq: publisher}
}

// NewRabbitSender creates a new rabbitMqSender instance base on connection and queue.
//
// Parameters:
//   - ctx: The context for the publisher.
//   - queueName: The name of the queue.
//   - conn: The connection to use for publishing messages, which stands for "rabbitmq.AmqpConnection".
//
// Returns:
//   - RabbitMqSender: An RabbitMqSender interface
//   - error
func NewRabbitSender(ctx context.Context, queueName string, conn *rmq.AmqpConnection) (RabbitMqSender, error) {
	publisher, err := conn.NewPublisher(ctx, &rmq.QueueAddress{Queue: queueName}, nil)
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}
	return NewRabbitSenderInstance(ctx, publisher), nil
}
