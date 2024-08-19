package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	"errors"
	// "io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
)

func JustMake(f APIFunc) http.Handler{
    fun := func(w http.ResponseWriter, r *http.Request) {
        if err := f(w,r); err != nil {
            if apiErr , ok := err.(HandlerError); ok {
                http.Error(w, apiErr.Error(),apiErr.StatusCode)
                return
            } else {
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        }
    }
    return http.HandlerFunc(fun)
}

type storageMock struct{
    val float64
    name string
    mtype string
}

func (s *storageMock)Add(name string, val interface{}) (error) {
    switch val := val.(type) {
    case metrics.Gauge:
        s.val = float64(val)
        s.mtype = "gauge"
        return nil
    case metrics.Counter:
        if s.mtype == "counter" {
            s.val += float64(val)
            return nil
        } else {
            s.val = float64(val)
            s.mtype = "counter"
        }
    }
    return errors.New("incorect metric type")
}

func (s *storageMock) Get(name string, val interface{}) (error) {
    v := reflect.ValueOf(val)
    v = v.Elem()
    if s.name == name {
        if v.CanFloat() {
            v.SetFloat(s.val)
        }
        if v.CanInt() {
            v.SetInt(int64(s.val))
        }
        return nil
    }
    return errors.New("error")
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
                response: "Error while updating\n",
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
    r.Method(http.MethodPost, "/update/gauge/{name}/{val}",JustMake(UpdateGauge(&storageMock{})))

    for _, test := range tests{
        t.Run(test.name, func(t *testing.T) {

            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()
            assert.Equal(t, test.want.code, res.StatusCode)

            // resBody, err := io.ReadAll(res.Body)
            //
            // require.NoError(t,err)
            //
            // assert.Equal(t,test.want.response, string(resBody))
            res.Body.Close()
        })
    }
}
func TestUpdateCounter(t *testing.T) {
    type want struct {
        code int
        val int64
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
                response: "Error while updating\n",
            },
        },
        {
            name: "All good",
            reqURL: "http://fuckintsite.com/update/counter/name/111",
            want: want{
                code: http.StatusOK,
                val:111,
                response: "",
            },
        },
        {
            name: "All good",
            reqURL: "http://fuckintsite.com/update/counter/name/111",
            want: want{
                code: http.StatusOK,
                val:222,
                response: "",
            },
        },
    }
    r := chi.NewRouter()
    stor := &storageMock{}
    r.Method(http.MethodPost, "/update/counter/{name}/{val}",JustMake(UpdateCounter(stor)))
    for _, test := range tests{
        t.Run(test.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()
            assert.Equal(t,test.want.code, res.StatusCode)

            // resBody, err := io.ReadAll(res.Body)
            //
            // require.NoError(t,err)
            //
            // assert.Equal(t,test.want.response, string(resBody))
            assert.Equal(t,test.want.val, int64(stor.val))
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
                response: "Metric not found\n",
            },
        },
    }
    stor := &storageMock{}

    r := chi.NewRouter()
    r.Method(http.MethodPost,"/value/gauge/{name}",JustMake(GetGauge(stor)))

    for _, test := range tests{
        stor.val = test.val
        stor.name = test.name
        t.Run(test.name, func(t *testing.T) {

            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()
            assert.Equal(t,test.want.code, res.StatusCode)

            // resBody, err := io.ReadAll(res.Body)
            //
            // require.NoError(t,err)
            //
            // assert.Equal(t,test.want.response, string(resBody))

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
                response: "Metric not found\n",
            },
        },
    }
    stor := &storageMock{}

    r := chi.NewRouter()
    r.Method(http.MethodPost, "/value/counter/{name}", JustMake(GetCounter(stor)))

    for _, test := range tests{
        stor.val = test.val
        stor.name = test.name
        t.Run(test.name, func(t *testing.T) {

            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()
            assert.Equal(t,test.want.code, res.StatusCode)
            //
            // resBody, err := io.ReadAll(res.Body)
            //
            // require.NoError(t,err)
            //
            // assert.Equal(t,test.want.response, string(resBody))

            res.Body.Close()
        })
    }
}
