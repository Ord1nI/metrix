package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"

	"github.com/Ord1nI/metrix/internal/handlers"
	"github.com/Ord1nI/metrix/internal/logger"
)

type sResponseWriter struct {
    http.ResponseWriter
    Signer *Signer
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

func MW(l logger.Logger, key []byte) func(http.Handler) http.Handler{
    return func (handler http.Handler) http.Handler {
        f := func(w http.ResponseWriter, r *http.Request) {
            stringHash := r.Header.Get("HashSHA256")
            if len(stringHash) < 64 {
                l.Infoln("Bad hash")
                w.WriteHeader(http.StatusBadRequest)
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

            signer := New(sha256.New, key)

            Hash, err := signer.Sign(bodyBytes)

            if err != nil {
                l.Infoln("Error while signing")
                handlers.SendInternalError(w)
                return
            }

            if !hmac.Equal(getHash, Hash) {
                l.Infoln("Hashes not equal")
                w.WriteHeader(http.StatusBadRequest)
                w.Write(nil)
                return
            }

            srw := &sResponseWriter{w, signer}

            l.Infoln("Request accepted")
            handler.ServeHTTP(srw, r)

            w.Header().Add("HashSHA256", hex.EncodeToString(srw.Signer.Sum(nil)))
        }
        return http.HandlerFunc(f)
    }
}
