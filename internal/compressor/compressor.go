package compressor

import (
    "compress/gzip"
    "strings"
    "io"
    "fmt"
    "net/http"
    "bytes"
    "errors"
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

func GzipMiddleware(h http.Handler) http.Handler{
    logFn := func(w http.ResponseWriter, r *http.Request) {
        if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
            newBody, err := NewGzipBody(r.Body)
            if err != nil {
                http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
            }
            r.Body = newBody
        }

        if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && 
            (strings.Contains(r.Header.Get("Content-Type"), "application/json") ||
                strings.Contains(r.Header.Get("Content-Type"), "text/html")) {
            gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)

            if err != nil {
                h.ServeHTTP(w,r) //be careful
            }

            defer gz.Close()

            w.Header().Set("Content-Encoding", "gzip")

            h.ServeHTTP(gzipWriter{w,gz}, r)
            return
        }

        h.ServeHTTP(w,r)

    }

    return http.HandlerFunc(logFn)
}
