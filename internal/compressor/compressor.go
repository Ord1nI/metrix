package compressor

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
)
type gzipWriter struct {
    http.ResponseWriter
    Writer io.Writer
}

type gzipBody struct {
    gz *gzip.Reader
    body io.ReadCloser
}

func NewGzipBody(body io.ReadCloser) (*gzipBody, error){
    gz, err := gzip.NewReader(body)

    if err != nil {
        return nil, err
    }
    return &gzipBody{
        gz:gz,
        body:body,
    }, nil
}

func (b *gzipBody) Read(p []byte) (n int, err error) {
    return b.gz.Read(p)
}

func (b *gzipBody) Close() error{
    err := errors.Join(
        b.gz.Close(),
        b.body.Close())
    return err
}

func (w gzipWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

func ToGzip(data []byte) ([]byte, error){
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

    return buf.Bytes(),nil
}

func FromGzip(data []byte)  ([]byte, error) {
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

