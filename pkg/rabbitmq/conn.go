package rabbitmq

import (
	"context"
	"fmt"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

// initQueue creates a new queue base on connection and queue name
func initQueue(ctx context.Context, conn *rmq.AmqpConnection, queueName string) error {
	_, err := conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{Name: queueName})
	if err != nil {
		conn.Close(ctx)
		return err
	}
	return nil
}

// NewConnection creates a new rabbitMqSender instance and create/get queue
//
// Parameters:
//   - ctx: The context for the publisher.
//   - queueName: The name of the queue.
//
// Returns:
//   - RabbitMqSender: An RabbitMqSender interface
func NewConnection(ctx context.Context, queueName string) (*rmq.AmqpConnection, error) {
	config, err := newConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	env := rmq.NewEnvironment(config.GetBrokerURI(), nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	err = initQueue(ctx, conn, queueName)
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to init queue: %w", err)
	}

	return conn, err

}
