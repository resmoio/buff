package buff

import (
	"sync"
	"time"
)

type Buffer[T any] struct {
	mu        sync.Mutex
	items     []T
	flushFunc func([]T)
	size      int
	timeout   time.Duration
	ticker    *time.Ticker
	done      chan struct{}
	closed    bool
}

func NewBuffer[T any](flushFunc func([]T), flushSize int, flushTimeout time.Duration) *Buffer[T] {
	buffer := &Buffer[T]{
		items:     make([]T, 0, flushSize),
		flushFunc: flushFunc,
		size:      flushSize,
		timeout:   flushTimeout,
		done:      make(chan struct{}),
	}
	return buffer
}

func (b *Buffer[T]) Push(item T) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		panic("tried to push on already closed buffer")
	}

	b.items = append(b.items, item)
	if len(b.items) >= b.size {
		b.flush()
	}
}

func (b *Buffer[T]) flush() {
	if len(b.items) == 0 {
		return
	}
	b.flushFunc(b.items)
	b.items = make([]T, 0, b.size)
	b.ticker.Reset(b.timeout)
}

func (b *Buffer[T]) Start() {
	b.ticker = time.NewTicker(b.timeout)
	go func() {
		for {
			select {
			case <-b.ticker.C:
				b.mu.Lock()
				b.flush()
				b.mu.Unlock()
			case <-b.done:
				b.ticker.Stop()
				return
			}
		}
	}()
}

func (b *Buffer[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.closed {
		b.closed = true
		close(b.done)
		b.flush()
	}
}
