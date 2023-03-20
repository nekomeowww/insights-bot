package utils

import (
	"context"
	"sync"

	"github.com/nekomeowww/insights-bot/pkg/options"
)

type InvokeOptions struct {
	ctx context.Context
}

func WithContext(ctx context.Context) options.CallOptions[InvokeOptions] {
	return options.NewCallOptions(func(o *InvokeOptions) {
		o.ctx = ctx
	})
}

func Invoke0(funcToBeRan func() error, callOpts ...options.CallOptions[InvokeOptions]) error {
	opts := options.ApplyCallOptions(callOpts, InvokeOptions{ctx: context.Background()})

	resChan := make(chan struct{}, 1)
	var err error

	go func() {
		err = funcToBeRan()
		resChan <- struct{}{}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case <-opts.ctx.Done():
			err = opts.ctx.Err()
		case <-resChan:
		}

		wg.Done()
	}()

	wg.Wait()
	return err
}
