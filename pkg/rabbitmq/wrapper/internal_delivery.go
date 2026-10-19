package wrapper

import (
	"context"

	"github.com/Azure/go-amqp"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

//go:generate mockery --name AmqpDelivery --filename delivery_mock.go --structname MockAmqpDelivery
type AmqpDelivery interface {
	Message() *amqp.Message
	Accept(ctx context.Context) error
}

// -- instance Internal AmqpDelivery
type realDeliveryWrapper struct {
	delivery rmq.IDeliveryContext // Đây là struct thật từ thư viện v1.4.0
}

// NewDeliveryWrapper init internal AmqpDelivery instance, implement for rmq.IDeliveryContext
func NewDeliveryWrapper(delivery rmq.IDeliveryContext) AmqpDelivery {
	return &realDeliveryWrapper{delivery: delivery}
}

// Delivery message
//
// Returns:
//   - *amqp.Message
func (w *realDeliveryWrapper) Message() *amqp.Message {
	return w.delivery.Message()
}

// Accept delivery
//
// Parameters:
//   - ctx: The context for the publisher.
//
// Returns:
//   - error
func (w *realDeliveryWrapper) Accept(ctx context.Context) error {
	return w.delivery.Accept(ctx)
}
