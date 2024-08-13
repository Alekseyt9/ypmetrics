package compress

import (
	"bytes"
	"compress/gzip"
	"io"
)

// GZIPCompress compresses the given data using gzip compression.
// It returns the compressed data as a byte slice or an error if the compression fails.
// Parameters:
//   - data: the data to be compressed
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

// GZIPDecompress decompresses the given gzip-compressed data.
// It returns the decompressed data as a byte slice or an error if the decompression fails.
// Parameters:
//   - compressedData: the data to be decompressed
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

// GZIPdecompressreader decompresses data read from the provided io.Reader using the given gzip.Reader.
// It returns the decompressed data as a byte slice or an error if the decompression fails.
// Parameters:
//   - reader: the io.Reader to read compressed data from
//   - gz: the gzip.Reader to use for decompression
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
