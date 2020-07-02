package main

import (
	"context"
	"testing"
	"time"
)

func assert(t *testing.T, ok bool, msg string, args ...interface{}) {
	if !ok {
		t.Errorf(msg, args...)
	}
}

func TestMutex_Lock(t *testing.T) {
	m := NewMutex()
	m.Lock()
	done := make(chan struct{})
	go func() {
		m.Lock()
		defer m.Unlock()
		close(done)
	}()
	m.Unlock()
	<-done
}

func TestMutex_TryLock(t *testing.T) {
	m := NewMutex()

	m.Lock()
	assert(t, !m.TryLock(), "TryLock must return false")
	m.Unlock()

	assert(t, m.TryLock(), "TryLock must return false")
	m.Unlock()
}

func TestMutex_TryLockWith(t *testing.T) {
	m := NewMutex()

	m.Lock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	assert(t, !m.TryLockWith(ctx), "TryLockWith must return false")
	m.Unlock()

	assert(t, m.TryLockWith(context.Background()), "TryLockWith must return true")
	m.Unlock()

	m.Lock()
	done := make(chan struct{})
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		ok := m.TryLockWith(ctx)
		assert(t, ok, "TryLockWith must return true")
		if ok {
			defer m.Unlock()
		}
		close(done)
	}()
	m.Unlock()
	<-done
}

func TestMutex_TryLockTimeout(t *testing.T) {
	m := NewMutex()

	m.Lock()
	assert(t, !m.TryLockTimeout(time.Second), "TryLockWith must return false")
	m.Unlock()

	assert(t, m.TryLockTimeout(time.Second), "TryLockWith must return true")
	m.Unlock()

	m.Lock()
	done := make(chan struct{})
	go func() {
		ok := m.TryLockTimeout(time.Second)
		assert(t, ok, "TryLockWith must return true")
		if ok {
			defer m.Unlock()
		}
		close(done)
	}()
	m.Unlock()
	<-done
}

func TestMutex_Unlock(t *testing.T) {
	m := NewMutex()
	defer func() {
		assert(t, recover() != nil, "unpaired Unlock must panic")
	}()
	m.Unlock()
}
