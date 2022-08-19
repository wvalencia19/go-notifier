package service

import (
	"context"
	"sync"
)

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
func totalChan(channel <-chan string) int {
	var wg sync.WaitGroup
	count := 0
	wg.Add(1)
	output := func(c <-chan string) {
		defer wg.Done()
		for range c {
			count++
		}
	}
	go output(channel)
	wg.Wait()
	return count
}

func cancelledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
