// Package pool provides pooling mechanisms for resources such as gzip readers.
package pool

import (
	"compress/gzip"
	"log/slog"
	"sync"
)

// ZipReaderPoolSingleton is a singleton structure that provides a pool of gzip readers.
type ZipReaderPoolSingleton struct {
	WriterPool sync.Pool
	log        *slog.Logger
}

var zrinstance *ZipReaderPoolSingleton
var zronce sync.Once

// GetZipReaderPool returns the singleton instance of ZipReaderPoolSingleton.
// It initializes the instance if it hasn't been initialized yet.
// Parameters:
//   - log: the logger to be used by the pool
func GetZipReaderPool(log *slog.Logger) *ZipReaderPoolSingleton {
	zronce.Do(func() {
		zrinstance = &ZipReaderPoolSingleton{
			WriterPool: sync.Pool{
				New: func() interface{} {
					reader := new(gzip.Reader)
					return reader
				},
			},
			log: log,
		}
	})
	return zrinstance
}
