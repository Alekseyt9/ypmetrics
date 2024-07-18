package pool

import (
	"compress/gzip"
	"sync"

	"github.com/Alekseyt9/ypmetrics/internal/server/log"
)

type ZipWriterPoolSingleton struct {
	WriterPool sync.Pool
	log        log.Logger
}

var zwinstance *ZipWriterPoolSingleton
var zwonce sync.Once

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
