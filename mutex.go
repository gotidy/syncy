package main

import (
	"context"
	"time"
)

// A Mutex is a mutual exclusion lock.
type Mutex struct {
	ch chan struct{}
}

// NewMutex create new Mutext.
func NewMutex() Mutex {
	return Mutex{ch: make(chan struct{}, 1)}
}

// Lock locks mutex
func (m Mutex) Lock() {
	m.ch <- struct{}{}
}

// TryLockWith try to lock the mutex if the context will not be done before
func (m Mutex) TryLock() bool {
	select {
	case m.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// TryLockWith trys to lock the mutex if the context will not be done before
func (m Mutex) TryLockCancel(cancel <-chan struct{}) bool {
	select {
	case m.ch <- struct{}{}:
		return true
	case <-cancel:
		return false
	}
}

// TryLockWith trys to lock the mutex if the context will not be done before
func (m Mutex) TryLockWith(ctx context.Context) bool {
	return m.TryLockCancel(ctx.Done())
}

// TryLock try to lock the mutex if timeout will not be happened before
func (m Mutex) TryLockTimeout(d time.Duration) bool {
	select {
	case m.ch <- struct{}{}:
		return true
	case <-time.After(d):
		return false
	}
}

// Unlock unlocks m. It is a run-time error if m is not locked on entry to Unlock.
func (m Mutex) Unlock() {
	select {
	case <-m.ch:
	default:
		panic("unlock is not paired with lock")
	}
}
