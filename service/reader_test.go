package service

import (
	"context"
	"io"
	mockReader "notifier/service/mocks" //nolint:typecheck
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader_Read(t *testing.T) {
	tests := []struct {
		name             string
		file             io.Reader
		ctx              context.Context
		messagesExpected int
		wantErr          bool
	}{
		{
			name:             "with correct data",
			file:             strings.NewReader("message1\nmessage2\nmessage3\nmessage4\nmessage5\nmessage6"),
			ctx:              context.Background(),
			messagesExpected: 6,
			wantErr:          false,
		},
		{
			name:             "with canceled context",
			file:             strings.NewReader("message1\nmessage2\nmessage2"),
			ctx:              cancelledContext(),
			messagesExpected: 0,
			wantErr:          false,
		},
		{
			name:             "with corrupted file",
			file:             mockReader.Reader("corrupted file"),
			ctx:              context.Background(),
			messagesExpected: 0,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.file, defaultBufferSize)
			messageChan, errChan := reader.Read(tt.ctx)
			messagesProcessed := totalChan(messageChan)
			resultError := getError(errChan)
			if resultError != nil && tt.wantErr == false {
				t.Errorf("unexpected error %v", resultError)
			}
			assert.Equal(t, tt.messagesExpected, messagesProcessed)
		})
	}
}
