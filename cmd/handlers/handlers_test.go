package handlers

import (
    "github.com/go-chi/chi/v5"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"
    "errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type storageMock struct{
    val float64
    name string
}

func (s *storageMock) GetGauge(name string) (float64, error) {
    if name == s.name {
        return s.val, nil
    }
    return 0, errors.New("error")
}
func (s *storageMock) AddGauge(name string, val float64) {
}
func (s *storageMock) GetCounter(name string) (int64, error) {
    if name == s.name {
        return int64(s.val), nil
    }
    return 0, errors.New("error")
}
func (s *storageMock) AddCounter(name string, val int64) {
}


func TestUpdateGauge(t *testing.T) {

    type want struct {
        code int
        response string
    }
    tests := []struct{
        name string
        reqURL string
        want want
    }{
        {
            name: "Test badReq",
            reqURL: "http://fuckintsite.com/update/gauge/name/afs",
            want: want{
                code: http.StatusBadRequest,
                response: "Incorect metric value\n",
            },
        },
        {
            name: "All good",
            reqURL: "http://fuckintsite.com/update/gauge/name/111.32",
            want: want{
                code: http.StatusOK,
                response: "",
            },
        },
    }

    r := chi.NewRouter()
    r.Post("/update/gauge/{name}/{val}",UpdateGauge(&storageMock{}))

    for _, test := range tests{
        t.Run(test.name, func(t *testing.T) {

            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()
            assert.Equal(t,test.want.code, res.StatusCode)

            resBody, err := io.ReadAll(res.Body)

            require.NoError(t,err)
            
            assert.Equal(t,test.want.response, string(resBody))
            res.Body.Close()
        })
    }
}
func TestUpdateCounter(t *testing.T) {
    type want struct {
        code int
        response string
    }
    tests := []struct{
        name string
        reqURL string
        want want
    }{
        {
            name: "Test Invorect value",
            reqURL: "http://fuckintsite.com/update/counter/name/123.34",
            want: want{
                code: http.StatusBadRequest,
                response: "Incorect metric value\n",
            },
        },
        {
            name: "All good",
            reqURL: "http://fuckintsite.com/update/counter/name/111",
            want: want{
                code: http.StatusOK,
                response: "",
            },
        },
    }
    r := chi.NewRouter()
    r.Post("/update/counter/{name}/{val}",UpdateCounter(&storageMock{}))
    for _, test := range tests{
        t.Run(test.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()
            assert.Equal(t,test.want.code, res.StatusCode)

            resBody, err := io.ReadAll(res.Body)

            require.NoError(t,err)
            
            assert.Equal(t,test.want.response, string(resBody))
            res.Body.Close()
        })
    }
}
func TestGetGauge(t *testing.T) {
    type want struct {
        code int
        response string
    }
    tests := []struct{
        reqURL string
        name string
        val float64
        want want
    }{
        {
            reqURL: "http://fuckintsite.com/value/gauge/name",
            name: "name",
            val: 123.324,
            want: want{
                code: http.StatusOK,
                response: "123.324\n",
            },
        },
        {
            reqURL: "http://fuckintsite.com/value/gauge/name12",
            name: "name123",
            want: want{
                code: http.StatusNotFound,
                response: "Unknown metric\n",
            },
        },
    }
    var stor *storageMock = &storageMock{}

    r := chi.NewRouter()
    r.Post("/value/gauge/{name}",GetGauge(stor))

    for _, test := range tests{
        stor.val = test.val
        stor.name = test.name
        t.Run(test.name, func(t *testing.T) {

            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()
            assert.Equal(t,test.want.code, res.StatusCode)

            resBody, err := io.ReadAll(res.Body)

            require.NoError(t,err)
            
            assert.Equal(t,test.want.response, string(resBody))

            res.Body.Close()
        })
    }
}

func TestGetCounter(t *testing.T) {
    type want struct {
        code int
        response string
    }
    tests := []struct{
        reqURL string
        name string
        val float64
        want want
    }{
        {
            reqURL: "http://fuckintsite.com/value/counter/name",
            name: "name",
            val: 123,
            want: want{
                code: http.StatusOK,
                response: "123\n",
            },
        },
        {
            reqURL: "http://fuckintsite.com/value/counter/name12",
            name: "name123",
            want: want{
                code: http.StatusNotFound,
                response: "Unknown metric\n",
            },
        },
    }
    var stor *storageMock = &storageMock{}

    r := chi.NewRouter()
    r.Post("/value/counter/{name}",GetCounter(stor))

    for _, test := range tests{
        stor.val = test.val
        stor.name = test.name
        t.Run(test.name, func(t *testing.T) {

            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()
            assert.Equal(t,test.want.code, res.StatusCode)

            resBody, err := io.ReadAll(res.Body)

            require.NoError(t,err)
            
            assert.Equal(t,test.want.response, string(resBody))

            res.Body.Close()
        })
    }
}
