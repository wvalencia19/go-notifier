package library

import (
	"context"
	"errors"
	senderMock "notifier/library/mocks"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestSender_Send(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		interval         time.Duration
		prepare          func(sender *senderMock.MockClient)
		messageChan      <-chan string
		messagesExpected int
		wantErr          bool
	}{
		{
			name:     "with correct data",
			ctx:      context.Background(),
			interval: 10 * time.Millisecond,
			prepare: func(sm *senderMock.MockClient) {
				wg := sync.WaitGroup{}
				wg.Add(5)
				sm.
					EXPECT().
					SendNotification(gomock.Any(), gomock.Any()).
					Do(func(ctx, message interface{}) {
						defer wg.Done()
					}).
					Times(5).
					Return(nil)
				go func() {
					wg.Wait()
				}()
			},
			messageChan: populatedChannel(5),
			wantErr:     false,
		},
		{
			name:     "with canceled context",
			ctx:      cancelledContext(),
			interval: 10 * time.Millisecond,
			prepare: func(sm *senderMock.MockClient) {
				sm.
					EXPECT().
					SendNotification(gomock.Any(), gomock.Any()).
					Times(0)
			},
			messageChan: populatedChannel(1),
			wantErr:     false,
		},
		{
			name:     "with error in http response",
			ctx:      context.Background(),
			interval: 10 * time.Millisecond,
			prepare: func(sm *senderMock.MockClient) {
				wg := sync.WaitGroup{}
				wg.Add(1)
				sm.
					EXPECT().
					SendNotification(gomock.Any(), gomock.Any()).
					Do(func(ctx, message interface{}) {
						defer wg.Done()
					}).
					Times(1).
					Return(errors.New("error performing http request"))
				go func() {
					wg.Wait()
				}()
			},
			messageChan: populatedChannel(1),
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sm := senderMock.NewMockClient(ctrl)
			tt.prepare(sm)
			sender := NewSender(tt.interval, sm)
			errList := sender.Send(tt.ctx, tt.messageChan)
			err := getError(errList)
			if err != nil && tt.wantErr == false {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
}

func populatedChannel(total int) <-chan string {
	messageChan := make(chan string, 10)
	defer close(messageChan)
	for i := 0; i < total; i++ {
		messageChan <- "message"
	}
	return messageChan
}

func getError(errList <-chan error) error {
	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	output := func(errList <-chan error) {
		defer wg.Done()
		for e := range errList {
			if e != nil {
				err = e
				return
			}
		}
	}
	go output(errList)
	wg.Wait()
	return err
}

func cancelledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
