package middlewares

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Ord1nI/metrix/internal/handlers"
)

type fileWriter interface {
	WriteToFile(f string) error
}

type logger interface {
	Errorln(args ...interface{})
	Infoln(args ...interface{})
}

func FileWriterWM(logger logger, stor fileWriter, path string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
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

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type gzipBody struct {
	gz   *gzip.Reader
	body io.ReadCloser
}

func newGzipBody(body io.ReadCloser) (*gzipBody, error) {
	gz, err := gzip.NewReader(body)

	if err != nil {
		return nil, err
	}
	return &gzipBody{
		gz:   gz,
		body: body,
	}, nil
}

func (b *gzipBody) Read(p []byte) (n int, err error) {
	return b.gz.Read(p)
}

func (b *gzipBody) Close() error {
	err := errors.Join(
		b.gz.Close(),
		b.body.Close())
	return err
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressorMW(l logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				newBody, err := newGzipBody(r.Body)
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

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		body         []byte
		responseData *responseData
	}
)

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	responseData := &responseData{
		status: 0,
		size:   0,
	}

	lw := loggingResponseWriter{
		ResponseWriter: w,
		responseData:   responseData,
	}
	return &lw
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	r.body = b
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggerMW(logger logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lw := newLoggingResponseWriter(w)

			h.ServeHTTP(lw, r)

			duration := time.Since(start)

			logger.Infoln(
				"\nREQUESE\n",
				"uri:", r.RequestURI, "\n",
				"method:", r.Method, "\n",
				"Header", r.Header, "\n",
				"RESPONSE\n",
				"status:", lw.responseData.status, "\n",
				"duration:", duration, "\n",
				"Header:", lw.ResponseWriter.Header(), "\n",
				"size:", lw.responseData.size,
			)
		}
		return http.HandlerFunc(logFn)
	}
}

type sResponseWriter struct {
    http.ResponseWriter
    Signer hash.Hash
}

type ReqBody struct {
    *bytes.Buffer
}

func (r *ReqBody) Close() error {
    r.Buffer.Reset()
    return nil
}

func (rw *sResponseWriter) Write(b []byte) (int,error) {
    n, err := rw.ResponseWriter.Write(b)
    _, err1 := rw.Signer.Write(b)
    return n, errors.Join(err,err1)
}


func SingMW(l logger, key []byte) func(http.Handler) http.Handler{
    return func (handler http.Handler) http.Handler {
        f := func(w http.ResponseWriter, r *http.Request) {
            stringHash := r.Header.Get("HashSHA256")
            if stringHash != "" {
                if len(stringHash) < 64 {
                    l.Infoln("Bad hash")
                    w.WriteHeader(http.StatusNotFound)
                    w.Write(nil)
                    return
                }

                getHash, err := hex.DecodeString(stringHash)
                if err != nil {
                    l.Infoln("error whiele decoding hex", err)
                    handlers.SendInternalError(w)
                    return
                }

                bodyBytes, err := io.ReadAll(r.Body)

                if err != nil {
                    l.Infoln("error while reading body", err)
                    handlers.SendInternalError(w)
                    return
                }

                defer r.Body.Close()
                r.Body = &ReqBody{
                    Buffer: bytes.NewBuffer(bodyBytes),
                }

                signer := hmac.New(sha256.New, key)
                _,err = signer.Write(bodyBytes)

                if err != nil {
                    l.Infoln("Error while signing")
                    handlers.SendInternalError(w)
                    return
                }

                Hash := signer.Sum(nil)

                if !hmac.Equal(getHash, Hash) {
                    l.Infoln("Hashes not equal")
                    l.Infoln(getHash, "\n", Hash)
                    w.WriteHeader(http.StatusBadRequest)
                    w.Write(nil)
                    return
                }

                signer.Reset()
                srw := &sResponseWriter{w, signer}

                l.Infoln("Request accepted")
                handler.ServeHTTP(srw, r)

                w.Header().Add("HashSHA256", hex.EncodeToString(srw.Signer.Sum(nil)))
            } else {
                handler.ServeHTTP(w,r)
            }
        }
        return http.HandlerFunc(f)
    }
}
