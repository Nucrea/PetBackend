package logger

import (
	"bufio"
	"context"
	"io"
	"sync"
	"time"
)

const FlushInterval = 500 * time.Millisecond

type bufioWrapper struct {
	writer *bufio.Writer
	ticker *time.Ticker
	mutex  *sync.RWMutex
}

func newWrapper(writer io.Writer) *bufioWrapper {
	ticker := time.NewTicker(FlushInterval)
	ticker.Stop()

	return &bufioWrapper{
		writer: bufio.NewWriterSize(writer, 128*1024),
		mutex:  &sync.RWMutex{},
		ticker: ticker,
	}
}

func (b *bufioWrapper) FlushRoutine(ctx context.Context) {
	go func() {
		b.ticker.Reset(FlushInterval)
		defer b.ticker.Stop()

		for {
			b.flush()

			select {
			case <-ctx.Done():
				b.flush()
				return
			case <-b.ticker.C:
			}
		}
	}()
}

func (b *bufioWrapper) Write(p []byte) (nn int, err error) {
	// TODO: try replace mutex, improve logging perfomance
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if len(p) > b.writer.Available() {
		b.ticker.Reset(FlushInterval)
	}

	return b.writer.Write(p)
}

func (b *bufioWrapper) flush() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.writer.Flush()
}

func (b *bufioWrapper) Close() error {
	return b.flush()
}
