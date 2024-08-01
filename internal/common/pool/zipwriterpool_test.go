package pool_test

import (
	"compress/gzip"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Alekseyt9/ypmetrics/internal/common/pool"
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
)

func TestGetZipWriterPool(t *testing.T) {
	mockLogger := log.NewNoOpLogger()

	zwp1 := pool.GetZipWriterPool(mockLogger)
	require.NotNil(t, zwp1)
	assert.NotNil(t, zwp1.WriterPool)

	zwp2 := pool.GetZipWriterPool(mockLogger)
	assert.Equal(t, zwp1, zwp2)
}

func TestZipWriterPoolSingleton(t *testing.T) {
	mockLogger := log.NewNoOpLogger()
	zwp := pool.GetZipWriterPool(mockLogger)

	var wg sync.WaitGroup
	const numGoroutines = 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			writer, ok := zwp.WriterPool.Get().(*gzip.Writer)
			require.True(t, ok)
			assert.NotNil(t, writer)
			zwp.WriterPool.Put(writer)
		}()
	}

	wg.Wait()
}
