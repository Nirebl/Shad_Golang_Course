//go:build !solution

package dupcall

import (
	"context"
	"sync"
)

type Call struct {
	mu sync.Mutex
	ch chan struct{}

	resources interface{}
	err       error
}

func (o *Call) gor(ctx context.Context, cb func(context.Context) (interface{}, error)) {
	o.resources, o.err = cb(ctx)
	o.mu.Lock()
	close(o.ch)
	o.ch = nil
	o.mu.Unlock()
}
func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.mu.Lock()
	if o.ch == nil {
		o.ch = make(chan struct{})
		go o.gor(ctx, cb)
	}

	lch := o.ch
	o.mu.Unlock()

	select {
	case <-lch:
		o.mu.Lock()
		defer o.mu.Unlock()
		return o.resources, o.err

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
