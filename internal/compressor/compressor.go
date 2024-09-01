package compressor

import (
	"bytes"
	"compress/gzip"
)

func ToGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)

	_, err := w.Write(data)

	if err != nil {
		return nil, err
	}

	err = w.Close()

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func FromGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	defer r.Close()

	var b bytes.Buffer

	_, err = b.ReadFrom(r)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
