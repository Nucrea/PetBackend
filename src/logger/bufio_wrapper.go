package logger

import (
	"bufio"
	"sync"
)

type BufioWrapper struct {
	*bufio.Writer
	m *sync.RWMutex
}

func (b *BufioWrapper) Write(p []byte) (nn int, err error) {
	b.m.RLock()
	defer b.m.RUnlock()

	return b.Writer.Write(p)
}

func (b *BufioWrapper) Flush() error {
	b.m.Lock()
	defer b.m.Unlock()

	return b.Writer.Flush()
}
