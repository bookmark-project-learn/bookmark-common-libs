package rabbitmq

import (
	"context"
	"fmt"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq/wrapper"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

//go:generate mockery --name RabbitMqReceiver --filename receiver.go
type RabbitMqReceiver interface {
	ConsumeOne(ctx context.Context) error
	Close() error
}

type rabbitMqReceiver struct {
	ctx         context.Context
	consumerRmq wrapper.InternalConsumer // Accept struct realConsumerWrapper và mockAmqpConsumer of test
}

// NewRabbitReceiverInstance init rabbitMqReceiver instance
func NewRabbitReceiverInstance(ctx context.Context, consumerRmq wrapper.InternalConsumer) RabbitMqReceiver {
	return &rabbitMqReceiver{
		ctx:         ctx,
		consumerRmq: consumerRmq,
	}
}

// NewRabbitReceiver creates a new rabbitMqReceiver instance base on connection and queue.
//
// Parameters:
//   - ctx: The context for the publisher.
//   - queueName: The name of the queue.
//   - conn: The connection to use for publishing messages, which stands for "rabbitmq.AmqpConnection".
//
// Returns:
//   - RabbitMqSender: An RabbitMqSender interface
//   - error
func NewRabbitReceiver(ctx context.Context, queueName string, conn *rmq.AmqpConnection) (RabbitMqReceiver, error) {
	consumer, err := conn.NewConsumer(ctx, queueName, &rmq.ConsumerOptions{InitialCredits: 1})
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// Bọc consumer thật vào realConsumerWrapper để nó thoả mãn interface InternalConsumer
	return NewRabbitReceiverInstance(ctx, wrapper.NewRabbitReceiverInstance(consumer)), nil
}
