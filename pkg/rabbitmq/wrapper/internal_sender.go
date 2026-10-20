package wrapper

import (
	"context"

	"github.com/Azure/go-amqp"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

// *rmq.Publisher implements Interface InternalPublisher
//
//go:generate mockery --name InternalPublisher --filename publisher_mock.go --structname MockAmqpPublisher
type InternalPublisher interface {
	Publish(ctx context.Context, msg *amqp.Message) (*rmq.PublishResult, error)
	Close(ctx context.Context) error
}
