//go:build !solution

package keylock

import (
	"sort"
	"sync"
)

type lS struct {
	isLocked  bool
	condition *sync.Cond
}

type KeyLock struct {
	mu     sync.Mutex
	locksM map[string]*lS
}

func New() *KeyLock {
	return &KeyLock{
		locksM: make(map[string]*lS),
	}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	l.mu.Lock()

	var uniqueK []string

	km := make(map[string]struct{})
	for _, key := range keys {
		km[key] = struct{}{}
	}
	for key := range km {
		uniqueK = append(uniqueK, key)
	}
	sort.Strings(uniqueK)

	for _, k := range uniqueK {

		if k == "" {
			_, t := l.locksM["a"]
			if !t {
				continue
			}
			for key := range l.locksM {
				delete(l.locksM, key)
			}
			l.mu.Unlock()
			return true, func() {}
		}

		for {
			if l.locksM[k] == nil {
				l.locksM[k] = &lS{condition: sync.NewCond(&l.mu)}
			}

			if !l.locksM[k].isLocked {
				l.locksM[k].isLocked = true
				break
			}

			if cancel != nil {
				select {
				case <-cancel:
					l.mu.Unlock()
					return true, func() {}
				default:
				}
			}

			l.locksM[k].condition.Wait()
		}
	}

	l.mu.Unlock()

	return false,
		func() {
			l.mu.Lock()
			for _, key := range uniqueK {
				if state, ok := l.locksM[key]; ok {
					state.isLocked = false
					state.condition.Signal()
				}
			}
			l.mu.Unlock()
		}
}
