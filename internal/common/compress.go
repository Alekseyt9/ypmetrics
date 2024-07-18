package common

import (
	"bytes"
	"compress/gzip"
	"io"
)

func GZIPCompress(data []byte) ([]byte, error) {
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

func GZIPdecompressreader(reader io.Reader, gz *gzip.Reader) ([]byte, error) {
	err := gz.Reset(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
