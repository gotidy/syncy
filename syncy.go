package main

import (
	"context"
)

type Locker struct {
	ch chan struct{}
}

func NewLocker() Locker {
	return Locker{ch: make(chan struct{}, 1)}
}

func (l Locker) Lock() {
	l.ch <- struct{}{}
}

func (l Locker) TryLockWith(ctx context.Context) bool {
	select {
	case l.ch <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

func (l Locker) TryLock() bool {
	select {
	case l.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l Locker) Unlock() {
	select {
	case <-l.ch:
	default:
		panic("unlock is not paired with lock")
	}
}
