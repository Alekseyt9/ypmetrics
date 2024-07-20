package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HashSHA256 generates a SHA-256 hash of the given data using the provided key.
// It returns the hexadecimal representation of the hash.
// Parameters:
//   - data: the data to be hashed
//   - key: the key to use for the HMAC
func HashSHA256(data []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
