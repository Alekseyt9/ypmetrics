package common

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/andybalholm/brotli"
)

func BrotliCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w := brotli.NewWriterLevel(&buf, brotli.BestCompression)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func BrotliDecompress(compressedData []byte) ([]byte, error) {
	b := bytes.NewReader(compressedData)
	r := brotli.NewReader(b)
	decompressedData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}

func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)

	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GZIPDecompress(compressedData []byte) ([]byte, error) {
	b := bytes.NewReader(compressedData)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	decompressedData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return decompressedData, nil
}
