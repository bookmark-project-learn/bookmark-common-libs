package rabbitmq

import (
	"context"
	"errors"
	"fmt"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq/wrapper"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

//go:generate mockery --name RabbitMqSender --filename sender.go
type RabbitMqSender interface {
	Publish(ctx context.Context, message []byte) error
	Close(ctx context.Context) error
}

type rabbitMqSender struct {
	publisherRmq wrapper.InternalPublisher // Thay đổi từ *rmq.Publisher thành interface
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
func NewRabbitSenderInstance(publisher wrapper.InternalPublisher) RabbitMqSender {
	return &rabbitMqSender{publisherRmq: publisher}
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
	return NewRabbitSenderInstance(publisher), nil
}

// Close closes amqpPublisher
//
// Returns:
//   - error
func (s *rabbitMqSender) Close(ctx context.Context) error {
	return s.publisherRmq.Close(ctx)
}

// Publish sends a message to the queue
//
// Parameters:
//   - ctx: The context for the publisher.
//   - message: The message to send.
//
// Returns:
//   - error
func (s *rabbitMqSender) Publish(ctx context.Context, message []byte) error {
	if len(message) == 0 {
		return errors.New("cannot publish an empty message payload")
	}

	msg := rmq.NewMessage(message)
	_, err := s.publisherRmq.Publish(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to publish message via amqp 1.0: %w", err)
	}
	return nil
}
