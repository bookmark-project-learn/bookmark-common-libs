package rabbitmq

import (
	"context"
	"errors"
	"fmt"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

// Close closes amqpPublisher
//
// Returns:
//   - error
func (s *rabbitMqSender) Close() error {
	return s.publisherRmq.Close(s.ctx)
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
