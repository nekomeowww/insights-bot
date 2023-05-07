package utils

import (
	"context"
	"sync"
)

func Invoke0(ctx context.Context, funcToBeRan func() error) error {
	var err error
	resChan := make(chan struct{}, 1)

	go func() {
		err = funcToBeRan()
		resChan <- struct{}{}
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-resChan:
		}

		wg.Done()
	}()
	wg.Wait()

	return err
}
