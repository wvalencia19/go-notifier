package service

import (
	"context"
	"io"
	senderMock "notifier/library/mocks"
	mockReader "notifier/service/mocks"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

const defaultInterval = 100 * time.Millisecond
const defaultBufferSize = 30

func TestProcessor_Process(t *testing.T) {
	tests := []struct {
		name        string
		file        io.Reader
		interval    time.Duration
		httpTimeout time.Duration
		maxErrors   int
		ctx         context.Context
		prepare     func(sender *senderMock.MockClient)
		wantErr     bool
	}{
		{
			name:     "first",
			file:     strings.NewReader("message1\nmessage2\nmessage3"),
			interval: defaultInterval,
			ctx:      context.Background(),
			prepare: func(sm *senderMock.MockClient) {
				wg := sync.WaitGroup{}
				wg.Add(3)
				sm.
					EXPECT().
					SendNotification(gomock.Any(), gomock.Any()).
					Do(func(ctx, message interface{}) {
						defer wg.Done()
					}).
					Times(3).
					Return(nil)
				go func() {
					wg.Wait()
				}()
			},
			wantErr: false,
		},
		{
			name:     "with corrupted file",
			file:     mockReader.Reader("corrupted file"),
			wantErr:  true,
			interval: defaultInterval,
			ctx:      context.Background(),
			prepare: func(sm *senderMock.MockClient) {
				sm.
					EXPECT().
					SendNotification(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:     "with canceled context",
			file:     strings.NewReader("message1\nmessage2\nmessage3"),
			wantErr:  false,
			interval: defaultInterval,
			ctx:      cancelledContext(),
			prepare: func(sm *senderMock.MockClient) {
				sm.
					EXPECT().
					SendNotification(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sm := senderMock.NewMockClient(ctrl)
			tt.prepare(sm)

			processor := NewProcessor(tt.file, sm, tt.interval, tt.maxErrors, defaultBufferSize)
			err := processor.Process(tt.ctx)
			resultError := getError(err)
			if resultError != nil && tt.wantErr == false {
				t.Errorf("unexpected error %v", resultError)
			}
		})
	}
}
