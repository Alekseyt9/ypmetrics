package common

import (
	"bytes"
	"testing"
)

func TestBrotliCompressDecompress(t *testing.T) {
	originalData := []byte("This is a test string for Brotli compression and decompression.")

	compressedData, err := BrotliCompress(originalData)
	if err != nil {
		t.Fatalf("BrotliCompress failed: %v", err)
	}

	decompressedData, err := BrotliDecompress(compressedData)
	if err != nil {
		t.Fatalf("BrotliDecompress failed: %v", err)
	}

	if !bytes.Equal(originalData, decompressedData) {
		t.Errorf("Decompressed data does not match original data. Got %s, want %s", decompressedData, originalData)
	}
}

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
