package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(data []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
