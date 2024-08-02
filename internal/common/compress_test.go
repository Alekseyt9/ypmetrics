package common

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGZIPCompressDecompress(t *testing.T) {
	originalData := []byte("This is a test string for GZIP compression and decompression.")

	compressedData, err := GZIPCompress(originalData)
	require.NoError(t, err, "GZIPCompress failed")

	decompressedData, err := GZIPDecompress(compressedData)
	require.NoError(t, err, "GZIPDecompress failed")

	assert.Equal(t, originalData, decompressedData, "Decompressed data does not match original data")
}

func TestGZIPdecompressreader(t *testing.T) {
	originalData := []byte("Hello, GZIP Compression!")

	compressedData, err := GZIPCompress(originalData)
	require.NoError(t, err, "Failed to compress data")

	gz, err := gzip.NewReader(bytes.NewReader(compressedData))
	require.NoError(t, err, "Failed to create gzip reader")
	defer gz.Close()

	decompressedData, err := GZIPdecompressreader(bytes.NewReader(compressedData), gz)
	require.NoError(t, err, "Failed to decompress data")

	assert.Equal(t, originalData, decompressedData, "Decompressed data does not match original")
}
