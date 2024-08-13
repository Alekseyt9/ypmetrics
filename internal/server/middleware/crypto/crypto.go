package crypto

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"

	"github.com/Alekseyt9/ypmetrics/internal/common/crypto"
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
)

type cryptoReadCloser struct {
	io.Reader
	body io.ReadCloser
}

func (crc *cryptoReadCloser) Close() error {
	return crc.body.Close()
}

func WithCrypto(next http.Handler, log log.Logger, pKey *rsa.PrivateKey) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {
		handleDecyper(w, r, next, log, pKey)
	}

	return http.HandlerFunc(compressFn)
}

func handleDecyper(w http.ResponseWriter, r *http.Request, next http.Handler, log log.Logger, pKey *rsa.PrivateKey) {
	encryptedData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Error reading request body", err)
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	var decryptedData []byte
	decryptedData, err = crypto.Decipher(encryptedData, pKey)

	if err != nil {
		log.Error("Error decrypting data", "error", err)
		http.Error(w, "Unable to decrypt request body", http.StatusInternalServerError)
		return
	}

	r.Body = &cryptoReadCloser{Reader: io.NopCloser(io.Reader(bytes.NewReader(decryptedData))), body: r.Body}
	next.ServeHTTP(w, r)
}
