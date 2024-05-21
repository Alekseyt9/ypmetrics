package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
)

type hashWriter struct {
	http.ResponseWriter
	Writer io.Writer
	Hash   hash.Hash
}

func (w hashWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	if err != nil {
		return n, err
	}

	w.Hash.Write(b)
	return n, nil
}

func WithHash(next http.Handler, hashKey string) http.Handler {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		if hashKey != "" {
			hash := hmac.New(sha256.New, []byte(hashKey))
			hw := &hashWriter{
				ResponseWriter: w,
				Writer:         w,
				Hash:           hash,
			}
			next.ServeHTTP(hw, r)

			hashSum := hw.Hash.Sum(nil)
			w.Header().Set("HashSHA256", hex.EncodeToString(hashSum))
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hashFn)
}
