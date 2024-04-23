package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w compressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithCompress(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			handleBr(w, r, next)
			return
		}
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handlegzip(w, r, next)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(logFn)
}

func handlegzip(w http.ResponseWriter, r *http.Request, next http.Handler) {
	wr, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer wr.Close()

	w.Header().Set("Content-Encoding", "gzip")
	next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: wr}, r)
}

func handleBr(w http.ResponseWriter, r *http.Request, next http.Handler) {
	var buf bytes.Buffer
	wr := brotli.NewWriterLevel(&buf, brotli.BestCompression)
	defer wr.Close()

	w.Header().Set("Content-Encoding", "br")
	next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: wr}, r)
}
