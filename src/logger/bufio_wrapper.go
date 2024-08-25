package logger

import (
	"bufio"
	"sync"
)

type bufioWrapper struct {
	*bufio.Writer
	m *sync.RWMutex
}

func (b *bufioWrapper) Write(p []byte) (nn int, err error) {
	// TODO: try replace mutex, improve logging perfomance
	b.m.RLock()
	defer b.m.RUnlock()

	return b.Writer.Write(p)
}

func (b *bufioWrapper) Flush() error {
	b.m.Lock()
	defer b.m.Unlock()

	return b.Writer.Flush()
}

func (b *bufioWrapper) Close() error {
	return b.Flush()
}
