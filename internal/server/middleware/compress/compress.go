package compress

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Alekseyt9/ypmetrics/internal/common/pool"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w compressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithCompress(next http.Handler, log *slog.Logger) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handlegzip(w, r, next, log)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(compressFn)
}

func handlegzip(w http.ResponseWriter, r *http.Request, next http.Handler, log *slog.Logger) {
	wps := pool.GetZipWriterPool(log)
	gz := wps.WriterPool.Get().(*gzip.Writer)
	defer wps.WriterPool.Put(gz)

	gz.Reset(w)
	defer gz.Close()

	w.Header().Set("Content-Encoding", "gzip")
	next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: gz}, r)
}
