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
	*gzip.Reader
	body io.ReadCloser
}

// Write compresses the data before writing it to the underlying ResponseWriter.
func (w compressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// WithCompress returns a middleware handler that adds gzip compression to the response if the client supports it.
func WithCompress(next http.Handler, log log.Logger) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handlegzip(w, r, next, log)
			return
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			handleDecompress(w, r, next, log)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(compressFn)
}

func handlegzip(w http.ResponseWriter, r *http.Request, next http.Handler, log log.Logger) {
	wps := pool.GetZipWriterPool(log)
	gz := wps.WriterPool.Get().(*gzip.Writer)
	defer wps.WriterPool.Put(gz)

	gz.Reset(w)
	defer gz.Close()

	w.Header().Set("Content-Encoding", "gzip")
	next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: gz}, r)
}

func handleDecompress(w http.ResponseWriter, r *http.Request, next http.Handler, log log.Logger) {
	wps := pool.GetZipReaderPool(log)
	gz := wps.WriterPool.Get().(*gzip.Reader)
	defer wps.WriterPool.Put(gz)
	r.Body = &gzipReadCloser{Reader: gz, body: r.Body}
	next.ServeHTTP(w, r)
}
