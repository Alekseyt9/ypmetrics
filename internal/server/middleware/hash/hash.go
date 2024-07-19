// Package hash provides middleware for an HTTP server that adds hash computation to responses.
// It uses HMAC with SHA-256 to compute hashes of the response bodies.
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

// Write writes bytes to the underlying Writer and updates the hash value with the data written.
func (w hashWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	if err != nil {
		return n, err
	}

	if w.Hash != nil {
		w.Hash.Write(b)
	}

	return n, nil
}

// WithHash returns a middleware handler that computes a SHA-256 hash of the response body
// if a hash key is provided. The resulting hash is set in the "HashSHA256"
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
