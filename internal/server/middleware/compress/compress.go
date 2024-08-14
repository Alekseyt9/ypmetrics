// Package compress provides middleware for an HTTP server that implements gzip compression for responses.
package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/Alekseyt9/ypmetrics/internal/common/pool"
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type gzipReadCloser struct {
	io.Reader
	body io.ReadCloser
}

func (crc *gzipReadCloser) Close() error {
	return crc.body.Close()
}

// Write compresses the data before writing it to the underlying ResponseWriter.
func (w compressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// WithCompress returns a middleware handler that adds gzip compression to the response if the client supports it.
func WithCompress(next http.Handler, log log.Logger) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {
		writer := w

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			wps := pool.GetZipWriterPool(log)
			gz := wps.WriterPool.Get().(*gzip.Writer)
			defer wps.WriterPool.Put(gz)
			gz.Reset(w)
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")
			writer = compressWriter{ResponseWriter: w, Writer: gz}
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			wps := pool.GetZipReaderPool(log)
			gz := wps.ReaderPool.Get().(*gzip.Reader)
			defer wps.ReaderPool.Put(gz)

			if err := gz.Reset(r.Body); err != nil {
				log.Error("Failed to reset gzip reader", "error", err)
				http.Error(w, "Invalid gzip body", http.StatusBadRequest)
				return
			}

			r.Body = &gzipReadCloser{Reader: gz, body: r.Body}
		}

		next.ServeHTTP(writer, r)
	}

	return http.HandlerFunc(compressFn)
}
