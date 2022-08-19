package service

import (
	"context"
	"io"
	"notifier/library"
	"sync"
	"time"
)

type Processor struct {
	file       io.Reader
	client     library.Client
	interval   time.Duration
	maxErrors  int
	bufferSize int
}

func NewProcessor(file io.Reader, client library.Client, interval time.Duration,
	maxErrors int, bufferSize int) Processor {
	return Processor{
		file: file, client: client, interval: interval, maxErrors: maxErrors, bufferSize: bufferSize,
	}
}

func (p *Processor) Process(ctx context.Context) <-chan error {
	var errList []<-chan error

	reader := NewReader(p.file, p.bufferSize)
	sender := library.NewSender(p.interval, p.client)

	messageChan, errChan := reader.Read(ctx)
	errList = append(errList, errChan)

	errChan = sender.Send(ctx, messageChan)
	errList = append(errList, errChan)

	errChan = p.mergeErrors(errList...)
	return errChan
}

func (p *Processor) mergeErrors(errList ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(errList))

	output := func(c <-chan error) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(errList))
	for _, c := range errList {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
