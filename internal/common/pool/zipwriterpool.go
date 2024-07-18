package pool

import (
	"compress/gzip"
	"log/slog"
	"sync"
)

type ZipWriterPoolSingleton struct {
	WriterPool sync.Pool
	log        *slog.Logger
}

var instance *ZipWriterPoolSingleton
var once sync.Once

func GetZipWriterPool(log *slog.Logger) *ZipWriterPoolSingleton {
	once.Do(func() {
		instance = &ZipWriterPoolSingleton{
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
	return instance
}
