package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w compressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithCompress(next http.Handler) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handlegzip(w, r, next)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(compressFn)
}

func handlegzip(w http.ResponseWriter, r *http.Request, next http.Handler) {
	wr, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		_, err1 := io.WriteString(w, err.Error())
		if err1 != nil {
			panic(err1)
		}
		return
	}
	defer wr.Close()

	w.Header().Set("Content-Encoding", "gzip")
	next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: wr}, r)
}
