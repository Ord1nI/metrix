package storage

import (
	"net/http"
    "strings"
)

func SaveToFileMW(path string, stor *MemStorage) func(http.Handler) http.Handler{
    return func(h http.Handler) http.Handler {
        f := func(w http.ResponseWriter, r *http.Request) {
            h.ServeHTTP(w,r)
            if strings.Contains(r.URL.String(), "update") {
                stor.WriteToFile(path) // add logger in future
            }
        }
        return http.HandlerFunc(f)
    }
}
