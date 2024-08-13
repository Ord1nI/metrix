package storage

import (
    "net/http"
)

func SaveToFileMW(path string, stor *MemStorage) func(http.Handler) http.Handler{
    return func(h http.Handler) http.Handler {
        f := func(w http.ResponseWriter, r *http.Request) {
            stor.WriteToFile(path) // add logger in future
            h.ServeHTTP(w,r)
        }
        return http.HandlerFunc(f)
    }
}
