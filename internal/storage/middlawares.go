package storage

import (
	"net/http"
    "strings"

    "github.com/Ord1nI/metrix/internal/logger"
)

func SaveToFileMW(logger logger.Logger, path string, stor *MemStorage) func(http.Handler) http.Handler{
    return func(h http.Handler) http.Handler {
        f := func(w http.ResponseWriter, r *http.Request) {
            h.ServeHTTP(w,r)
            if strings.Contains(r.URL.String(), "update") {
                err := stor.WriteToFile(path) // add logger in future
                if err != nil {
                    logger.Errorln("Error while wiring to file:", path)
                } else {
                    logger.Infoln("all data Successfuly loaded to file")
                }
            }
        }
        return http.HandlerFunc(f)
    }
}
