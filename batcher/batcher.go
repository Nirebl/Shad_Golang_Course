//go:build !solution

package batcher

import (
	"sync"
	"sync/atomic"
	"time"

	"gitlab.com/slon/shad-go/batcher/slow"
)

type Batcher struct {
	mu         sync.Mutex
	cond       *sync.Cond
	value      *slow.Value
	lastChange uint64
	time       time.Time
}

func NewBatcher(v *slow.Value) *Batcher {
	b := &Batcher{
		value: v,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}
func (b *Batcher) Load() interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	enter := time.Now()

	for enter.Before(time.Unix(0, int64(atomic.LoadUint64(&b.lastChange)))) {
		time.Sleep(1 * time.Millisecond)
	}

	value := b.value.Load()

	// Обновляем lastChange атомарно
	atomic.StoreUint64(&b.lastChange, uint64(time.Now().UnixNano()))

	return value
}

func (b *Batcher) Store(value interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.value.Store(value)
	atomic.StoreUint64(&b.lastChange, uint64(time.Now().UnixNano()))
}
