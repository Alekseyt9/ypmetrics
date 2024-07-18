package pool

import (
	"compress/gzip"
	"log/slog"
	"sync"
)

type ZipReaderPoolSingleton struct {
	WriterPool sync.Pool
	log        *slog.Logger
}

var zrinstance *ZipReaderPoolSingleton
var zronce sync.Once

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
