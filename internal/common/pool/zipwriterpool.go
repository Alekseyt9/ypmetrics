// Package pool provides pooling mechanisms for resources such as gzip writers.
package pool

import (
	"compress/gzip"
	"sync"

	"github.com/Alekseyt9/ypmetrics/internal/server/log"
)

// ZipWriterPoolSingleton is a singleton structure that provides a pool of gzip writers.
type ZipWriterPoolSingleton struct {
	WriterPool sync.Pool
	log        log.Logger
}

var zwinstance *ZipWriterPoolSingleton
var zwonce sync.Once

// GetZipWriterPool returns the singleton instance of ZipWriterPoolSingleton.
// It initializes the instance if it hasn't been initialized yet.
// Parameters:
//   - log: the logger to be used by the pool
func GetZipWriterPool(log log.Logger) *ZipWriterPoolSingleton {
	zwonce.Do(func() {
		zwinstance = &ZipWriterPoolSingleton{
			WriterPool: sync.Pool{
				New: func() interface{} {
					writer, err := gzip.NewWriterLevel(nil, gzip.BestSpeed)
					if err != nil {
						log.Error("error create gzip.NewWriterLevel", "error", err)
					}
					return writer
				},
			},
			log: log,
		}
	})
	return zwinstance
}
