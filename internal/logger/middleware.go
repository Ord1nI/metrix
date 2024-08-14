package logger

import (
    "go.uber.org/zap"

    "net/http"
    "time"
)

func HandlerLogging(logger *zap.SugaredLogger) func(http.Handler) http.Handler{
    return func(h http.Handler) http.Handler{
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
                "Header:",lw.ResponseWriter.Header(), "\n",
                "size:", lw.responseData.size,
            )
        }
        return http.HandlerFunc(logFn)
    }
}
