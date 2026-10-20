package test

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq/wrapper/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRabbitMqReceiver_ConsumeOne(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mConn *mocks.MockInternalConsumer, mDeliv *mocks.MockAmqpDelivery, ctx context.Context)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Nhận và xử lý tin nhắn thành công (Dữ liệu đơn lẻ)",
			setup: func(mConn *mocks.MockInternalConsumer, mDeliv *mocks.MockAmqpDelivery, ctx context.Context) {
				// 1. Giả lập gói tin AMQP 1.0 trả về từ Broker
				amqpMsg := &amqp.Message{
					Data: [][]byte{[]byte("hello-rabbitmq-4")},
				}
				mDeliv.On("Message").Return(amqpMsg)
				mDeliv.On("Accept", ctx).Return(nil)

				// 2. Cấu hình consumer trả về mock delivery trên
				mConn.On("Receive", ctx).Return(mDeliv, nil)
			},
			wantErr: false,
		},
		{
			name: "Nhận và xử lý tin nhắn thành công (Dữ liệu phân mảnh nhiều phần)",
			setup: func(mConn *mocks.MockInternalConsumer, mDeliv *mocks.MockAmqpDelivery, ctx context.Context) {
				// Giả lập tin nhắn lớn bị cắt nhỏ thành nhiều chunk
				amqpMsg := &amqp.Message{
					Data: [][]byte{[]byte("part-1-"), []byte("part-2")},
				}
				mDeliv.On("Message").Return(amqpMsg)
				mDeliv.On("Accept", ctx).Return(nil)

				mConn.On("Receive", ctx).Return(mDeliv, nil)
			},
			wantErr: false,
		},
		{
			name: "Lỗi đứt kết nối mạng khi Receive tin nhắn",
			setup: func(mConn *mocks.MockInternalConsumer, mDeliv *mocks.MockAmqpDelivery, ctx context.Context) {
				// Hàm Receive trả về lỗi kết nối thẳng từ thư viện gốc
				mConn.On("Receive", ctx).Return(nil, errors.New("amqp link detached"))
			},
			wantErr: true,
			errMsg:  "error receiving from amqp 1.0 link: amqp link detached",
		},
		{
			name: "Lỗi không thể Ack (Accept) tin nhắn về cho RabbitMQ",
			setup: func(mConn *mocks.MockInternalConsumer, mDeliv *mocks.MockAmqpDelivery, ctx context.Context) {
				amqpMsg := &amqp.Message{
					Data: [][]byte{[]byte("some-payload")},
				}
				mDeliv.On("Message").Return(amqpMsg)
				// Giả lập lỗi khi gửi lệnh Ack
				mDeliv.On("Accept", ctx).Return(errors.New("session closed"))

				mConn.On("Receive", ctx).Return(mDeliv, nil)
			},
			wantErr: true,
			errMsg:  "failed to accept delivery: session closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Khởi tạo các Mock object được sinh ra từ Mockery
			mockConsumer := mocks.NewMockInternalConsumer(t)
			mockDelivery := mocks.NewMockAmqpDelivery(t)
			ctx := context.Background()

			// Cấu hình kịch bản mock cho từng testcase
			tt.setup(mockConsumer, mockDelivery, ctx)

			// Tạo instance của receiver, inject mock consumer vào trực tiếp field internal
			receiver := rabbitmq.NewRabbitReceiverInstance(mockConsumer)
			// Chạy hàm kiểm thử ConsumeOne
			err := receiver.ConsumeOne(ctx)

			// Đánh giá kết quả bằng thư viện testify/assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
