package common

import (
	"bytes"
	"testing"
)

func TestGZIPCompressDecompress(t *testing.T) {
	originalData := []byte("This is a test string for GZIP compression and decompression.")

	compressedData, err := GZIPCompress(originalData)
	if err != nil {
		t.Fatalf("GZIPCompress failed: %v", err)
	}

	decompressedData, err := GZIPDecompress(compressedData)
	if err != nil {
		t.Fatalf("GZIPDecompress failed: %v", err)
	}

	if !bytes.Equal(originalData, decompressedData) {
		t.Errorf("Decompressed data does not match original data. Got %s, want %s", decompressedData, originalData)
	}
}
