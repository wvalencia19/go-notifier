package library

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Sender interface {
	Send(ctx context.Context, messageChan <-chan string) <-chan error
}

type sender struct {
	interval time.Duration
	client   Client
}

func NewSender(interval time.Duration, client Client) Sender {
	return &sender{interval: interval, client: client}
}

func (s *sender) Send(ctx context.Context, messageChan <-chan string) <-chan error {
	var wg sync.WaitGroup
	wg.Add(1)
	errChan := make(chan error, 1)

	go func() {
		defer wg.Done()
		interval := time.NewTicker(s.interval)
		defer interval.Stop()

		for message := range messageChan {
			select {
			case <-ctx.Done():
				log.Debug("return from reader due canceled context")
				return
			default:
			}

			select {
			case <-interval.C:
				{
					wg.Add(1)
					go func(m string) {
						defer wg.Done()
						err := s.client.SendNotification(ctx, m)
						if err != nil {
							errChan <- err
						}
					}(message)
				}
			case <-ctx.Done():
				log.Debug("return from sender due canceled context")
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()
	return errChan
}
