package compressor

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

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
            (strings.Contains(r.Header.Get("Accept"), "application/json") ||
                strings.Contains(r.Header.Get("Accept"), "html")) {

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
