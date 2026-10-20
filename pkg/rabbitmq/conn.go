package rabbitmq

import (
	"context"
	"fmt"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type RabbitClient struct {
	conn      *rmq.AmqpConnection
	queueName string
	publisher RabbitMqSender
	sender    RabbitMqSender
}

func buildConnection(ctx context.Context, queueName string) (*rmq.AmqpConnection, error) {
	config, err := newConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	env := rmq.NewEnvironment(config.GetBrokerURI(), nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	_, err = conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{Name: queueName})
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to init queue: %w", err)
	}
	return conn, nil
}

func InitInstance(conn *rmq.AmqpConnection, queueName string) *RabbitClient {
	return &RabbitClient{conn: conn, queueName: queueName}
}

// NewConnection creates a new rabbitMqSender instance and create/get queue
//
// Parameters:
//   - ctx: The context for the publisher.
//   - queueName: The name of the queue.
//
// Returns:
//   - RabbitMqSender: An RabbitMqSender interface
func NewConnection(ctx context.Context, queueName string) (*RabbitClient, error) {
	conn, err := buildConnection(ctx, queueName)
	if err != nil {
		return nil, err
	}
	return InitInstance(conn, queueName), nil
}

func (r *RabbitClient) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func (r *RabbitClient) GetSender(ctx context.Context) (RabbitMqSender, error) {
	if r.sender == nil {
		sender, err := NewRabbitSender(ctx, r.queueName, r.conn)
		if err != nil {
			return nil, err
		}
		r.sender = sender
	}
	return r.sender, nil
}

func (r *RabbitClient) GetPublisher(ctx context.Context) (RabbitMqSender, error) {
	if r.publisher == nil {
		publisher, err := NewRabbitSender(ctx, r.queueName, r.conn)
		if err != nil {
			return nil, err
		}
		r.publisher = publisher
	}
	return r.publisher, nil
}
