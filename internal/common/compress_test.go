package common

import (
	"bytes"
	"compress/gzip"
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

func TestGZIPdecompressreader(t *testing.T) {
	originalData := []byte("Hello, GZIP Compression!")

	compressedData, err := GZIPCompress(originalData)
	if err != nil {
		t.Fatalf("Failed to compress data: %v", err)
	}

	gz, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gz.Close()

	decompressedData, err := GZIPdecompressreader(bytes.NewReader(compressedData), gz)
	if err != nil {
		t.Fatalf("Failed to decompress data: %v", err)
	}

	if !bytes.Equal(decompressedData, originalData) {
		t.Errorf("Decompressed data does not match original. Got %v, want %v", decompressedData, originalData)
	}
}
