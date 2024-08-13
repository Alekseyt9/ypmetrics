package hash_test

import (
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/common/hash"
)

func TestHashSHA256(t *testing.T) {
	data := []byte("test")
	key := []byte("secret_key")
	expectedHash := "3c385748b9c2960d12944cf55e5bc9406f5ba79c2b942971a89c890c0b1f3a61"

	hash := hash.HashSHA256(data, key)

	if hash != expectedHash {
		t.Errorf("HashSHA256 failed: expected %s, got %s", expectedHash, hash)
	}
}
