package logger

import  (
    "go.uber.org/zap"

    "net/http"
)

func NewLogger() (*zap.Logger, error){
    logger, err := zap.NewDevelopment()
    return logger, err
}

type (

    responseData struct {
        status int
        size int
    }

    loggingResponseWriter struct {
        http.ResponseWriter
        responseData *responseData
    }
)

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
    responseData := &responseData {
        status: 0,
        size: 0,
    }

    lw := loggingResponseWriter {
        ResponseWriter: w,
        responseData: responseData,
    }
    return &lw
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
    size, err := r.ResponseWriter.Write(b)
    r.responseData.size += size
    return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
    r.ResponseWriter.WriteHeader(statusCode)
    r.responseData.status = statusCode
}

