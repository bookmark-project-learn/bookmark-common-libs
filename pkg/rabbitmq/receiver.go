package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq/wrapper"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

//go:generate mockery --name RabbitMqReceiver --filename receiver.go
type RabbitMqReceiver interface {
	ConsumeOne(ctx context.Context) error
	Close(ctx context.Context) error
}

type rabbitMqReceiver struct {
	consumerRmq wrapper.InternalConsumer // Accept struct realConsumerWrapper và mockAmqpConsumer of test
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
	return NewRabbitReceiverInstance(wrapper.NewRabbitReceiverInstance(consumer)), nil
}

// NewRabbitReceiverInstance init rabbitMqReceiver instance
func NewRabbitReceiverInstance(consumerRmq wrapper.InternalConsumer) RabbitMqReceiver {
	return &rabbitMqReceiver{
		consumerRmq: consumerRmq,
	}
}

// ConsumeOne consume one message
//
// Parameters:
//   - ctx: The context for the publisher.
//
// Returns:
//   - error
func (r *rabbitMqReceiver) ConsumeOne(ctx context.Context) error {
	delivery, err := r.consumerRmq.Receive(ctx)
	if err != nil {
		return fmt.Errorf("error receiving from amqp 1.0 link: %w", err)
	}

	msg := delivery.Message()
	var body string

	if msg != nil && len(msg.Data) > 0 {
		var completePayload []byte
		for _, chunk := range msg.Data {
			completePayload = append(completePayload, chunk...)
		}
		body = string(completePayload)
	}

	log.Printf("[AMQP 1.0] Received message body: %s", body)

	err = delivery.Accept(ctx)
	if err != nil {
		return fmt.Errorf("failed to accept delivery: %w", err)
	}

	return nil
}

// Close closes amqpPublisher
//
// Returns:
//   - error
func (r *rabbitMqReceiver) Close(ctx context.Context) error {
	return r.consumerRmq.Close(ctx)
}
