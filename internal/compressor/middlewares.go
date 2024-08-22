package compressor

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ord1nI/metrix/internal/logger"
)

func MW(l logger.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				newBody, err := NewGzipBody(r.Body)
				if err != nil {
					l.Errorln("req", r.URL.String, "Error while decoding gzip")
					http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
					return
				}
				r.Body = newBody
			}

			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") &&
				(strings.Contains(r.Header.Get("Accept"), "application/json") ||
					strings.Contains(r.Header.Get("Accept"), "html")) {

				gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)

				if err != nil {
					l.Errorln("Req", r.URL.String(), "Failed creating gzip encoder sending response without compression")
					h.ServeHTTP(w, r) //be careful
					return
				}

				defer gz.Close()

				w.Header().Set("Content-Encoding", "gzip")

				h.ServeHTTP(gzipWriter{w, gz}, r)
				return
			}

			h.ServeHTTP(w, r)

		}

		return http.HandlerFunc(fn)
	}
}
