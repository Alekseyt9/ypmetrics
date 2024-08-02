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

func TestGetZipReaderPool(t *testing.T) {
	mockLogger := log.NewNoOpLogger()

	zrp1 := pool.GetZipReaderPool(mockLogger)
	require.NotNil(t, zrp1)
	//assert.NotNil(t, zrp1.WriterPool)

	zrp2 := pool.GetZipReaderPool(mockLogger)
	assert.Equal(t, zrp1, zrp2)
}

func TestZipReaderPoolSingleton(t *testing.T) {
	mockLogger := log.NewNoOpLogger()
	zrp := pool.GetZipReaderPool(mockLogger)

	var wg sync.WaitGroup
	const numGoroutines = 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reader, ok := zrp.WriterPool.Get().(*gzip.Reader)
			require.True(t, ok)
			assert.NotNil(t, reader)
			zrp.WriterPool.Put(reader)
		}()
	}

	wg.Wait()
}
