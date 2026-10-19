package test

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

// 1. Mock struct sử dụng chính xác kiểu dữ liệu đã hiệu chỉnh
type mockAmqpPublisher struct {
	OnPublish func(ctx context.Context, msg *amqp.Message) (*rmq.PublishResult, error)
	OnClose   func(ctx context.Context) error
}

func (m *mockAmqpPublisher) Publish(ctx context.Context, msg *amqp.Message) (*rmq.PublishResult, error) {
	if m.OnPublish != nil {
		return m.OnPublish(ctx, msg)
	}
	return nil, nil // Trả về PublishResult trống (vì rmq.PublishResult là một interface hoặc kiểu con)
}

func (m *mockAmqpPublisher) Close(ctx context.Context) error {
	if m.OnClose != nil {
		return m.OnClose(ctx)
	}
	return nil
}

// 2. Viết Table-driven tests
func TestRabbitMqSender_Publish(t *testing.T) {
	tests := []struct {
		name    string
		message []byte
		mockFn  func(ctx context.Context, msg *amqp.Message) (*rmq.PublishResult, error)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Gửi tin nhắn thành công",
			message: []byte("hello rabbitmq"),
			mockFn: func(ctx context.Context, msg *amqp.Message) (*rmq.PublishResult, error) {
				if len(msg.Data) == 0 || string(msg.Data[0]) != "hello rabbitmq" {
					t.Errorf("Dữ liệu đóng gói message không chính xác")
				}
				return nil, nil
			},
			wantErr: false,
		},
		{
			name:    "Lỗi kết nối từ Broker",
			message: []byte("test error"),
			mockFn: func(ctx context.Context, msg *amqp.Message) (*rmq.PublishResult, error) {
				return nil, errors.New("network timeout")
			},
			wantErr: true,
			errMsg:  "failed to publish message via amqp 1.0: network timeout",
		},
		{
			name:    "Validation payload rỗng",
			message: []byte(""),
			mockFn: func(ctx context.Context, msg *amqp.Message) (*rmq.PublishResult, error) {
				return nil, nil
			},
			wantErr: true,
			errMsg:  "cannot publish an empty message payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPub := &mockAmqpPublisher{
				OnPublish: tt.mockFn,
			}
			ctx := context.Background()
			sender := rabbitmq.NewRabbitSenderInstance(ctx, mockPub)

			err := sender.Publish(ctx, tt.message)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Publish() error message = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}
