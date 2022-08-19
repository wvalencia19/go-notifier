package service

import (
	"bufio"
	"context"
	"io"

	log "github.com/sirupsen/logrus"
)

type Reader interface {
	Read(ctx context.Context) (<-chan string, <-chan error)
}

type reader struct {
	file       io.Reader
	bufferSize int
}

func NewReader(file io.Reader, bufferSize int) Reader {
	return reader{file: file, bufferSize: bufferSize}
}

func (r reader) Read(ctx context.Context) (<-chan string, <-chan error) {
	messageChan := make(chan string, r.bufferSize)
	errChan := make(chan error, 1)

	go func() {
		defer close(messageChan)
		defer close(errChan)
		scanner := bufio.NewScanner(r.file)
		log.Debugf("reading from file")

		for scanner.Scan() {
			b := scanner.Bytes()
			log.Debugf("reading line with value %s", string(b))

			select {
			case <-ctx.Done():
				log.Debug("return from reader due canceled context")
				return
			default:
			}

			select {
			case messageChan <- string(b):
			case <-ctx.Done():
				log.Debug("return from reader due canceled context")
				return
			}
		}
		if err := scanner.Err(); err != nil {
			errChan <- err
		}
	}()
	return messageChan, errChan
}
