package rabbitmq

import (
	"context"
	"fmt"
	"log"
)

func (r *rabbitMqReceiver) Close() error {
	return r.consumerRmq.Close(r.ctx)
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
