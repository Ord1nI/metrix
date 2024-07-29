package main

import (
    "io"
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/Ord1nI/metrix/cmd/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type want struct {
    code int
    res string
    val float64
}

type storageMock struct{
    val float64
}

func newSM() *storageMock{
    return &storageMock{
        val:0,
    }
}

func (s *storageMock) AddGauge(name string, val float64) {
    s.val = val
}

func (s *storageMock) AddCounter(name string, val int64) {
    s.val = float64(val)
}

func TestMain(t *testing.T) {
    stor := newSM()

    mux := http.NewServeMux()
    mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)})
    mux.HandleFunc(`/update/`, func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusBadRequest)})
    mux.HandleFunc(`/update/gauge/`, handlers.UpdateGauge(stor))
    mux.HandleFunc(`/update/counter/`, handlers.UpdateCounter(stor))

    serv := httptest.NewServer(mux)
    client := serv.Client()


    tests := []struct{
        URL string
        want want
    }{
        {
            URL: "",
            want: want{
                code:http.StatusNotFound,
                res: "",
            },
        },
        {
            URL: "/update",
            want: want{
                code:http.StatusBadRequest,
                res: "",
            },
        },
        {
            URL: "/update/gauge/",
            want: want{
                code:http.StatusNotFound,
                res: "Not Found\n",
            },
        },
        {
            URL: "/update/gauge/name/123/asdf/",
            want: want{
                code:http.StatusBadRequest,
                res: "Bad Request\n",
            },
        },
        {
            URL: "/update/gauge/name/ads",
            want: want{
                code:http.StatusBadRequest,
                res: "Incorect metric value\n",
            },
        },
        {
            URL: "/update/gauge/name/123.213",
            want: want{
                code:http.StatusOK,
                res: "",
                val: 123.213,
            },
        },
        {
            URL: "/update/gauge/name1/-123.213",
            want: want{
                code:http.StatusOK,
                res: "",
                val: -123.213,
            },
        },
    }
    for _, test := range tests {
        t.Run(serv.URL+test.URL,func(t *testing.T) {
            stor.val = 0

            res, err := client.Post(serv.URL+test.URL,"text/plain",nil)

            require.NoError(t,err)

            assert.Equal(t, test.want.code, res.StatusCode)

            r, _ := io.ReadAll(res.Body)
            res.Body.Close()
            assert.Equal(t, test.want.res, string(r))
            if stor == nil {
                stor.val = 0
            }
            assert.Equal(t, stor.val, test.want.val)
        })
    }
    
    tests = []struct{
        URL string
        want want
    }{
        {
            URL: "",
            want: want{
                code:http.StatusNotFound,
                res: "",
            },
        },
        {
            URL: "/update",
            want: want{
                code:http.StatusBadRequest,
                res: "",
            },
        },
        {
            URL: "/update/counter/",
            want: want{
                code:http.StatusNotFound,
                res: "Not Found\n",
            },
        },
        {
            URL: "/update/counter/name/123/asdf/",
            want: want{
                code:http.StatusBadRequest,
                res: "Bad Request\n",
            },
        },
        {
            URL: "/update/counter/name/ads",
            want: want{
                code:http.StatusBadRequest,
                res: "Incorect metric value\n",
            },
        },
        {
            URL: "/update/counter/name/123.213",
            want: want{
                code:http.StatusBadRequest,
                res: "Incorect metric value\n",
                val: 0,
            },
        },
        {
            URL: "/update/counter/name1/-123",
            want: want{
                code:http.StatusOK,
                res: "",
                val: -123,
            },
        },
        {
            URL: "/update/counter/name1/123",
            want: want{
                code:http.StatusOK,
                res: "",
                val: 123,
            },
        },
    }
    for _, test := range tests {
        t.Run(serv.URL+test.URL,func(t *testing.T) {
            stor.val = 0

            res, err := client.Post(serv.URL+test.URL,"text/plain",nil)

            require.NoError(t,err)

            assert.Equal(t, test.want.code, res.StatusCode)

            r, _ := io.ReadAll(res.Body)
            res.Body.Close()
            assert.Equal(t, test.want.res, string(r))
            if stor == nil {
                stor.val = 0
            }
            assert.Equal(t, stor.val, test.want.val)
        })
    }
}
