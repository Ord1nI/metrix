package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
    "github.com/go-chi/chi/v5"

    "io"
    "errors"
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/handlers"
)


type storageMock struct{
    val storage.Gauge
    name string
}

func newSM() *storageMock{
    return &storageMock{
        val:0,
    }
}

func (s *storageMock) AddGauge(name string, val storage.Gauge) {
    s.val = val
}

func (s *storageMock) AddCounter(name string, val storage.Counter) {
    s.val = storage.Gauge(val)
}

func (s *storageMock) GetGauge(name string) (storage.Gauge, error){
    if s.name == name {
        return s.val, nil
    }
    return 0, errors.New("error")
}

func (s *storageMock) GetCounter(name string) (storage.Counter, error){
    if s.name == name {
        return storage.Counter(s.val), nil
    }
    return 0, errors.New("error")
}

func TestMain(t *testing.T) {
    stor := newSM()

    r := chi.NewRouter()

    r.Route("/update", func(r chi.Router) {
        r.HandleFunc("/*", handlers.BadRequest)                      // ANY /update/


        r.Route("/gauge", updateGaugeRoute(stor))         // ANY /update/gauge/*

        r.Route("/counter", updateCounterRoute(stor))     // Any /update/counter/*
        
    })

    r.Route("/value", func(r chi.Router) {
        r.HandleFunc("/*", handlers.BadRequest)            // Any /value/

        r.Route("/gauge", valueGaugeRoute(stor))        

        r.Route("/counter", valueCounterRoute(stor))   
    })

    serv := httptest.NewServer(r)
    client := serv.Client()

    tCounter(t, stor, serv, client)
    tGauge(t, stor, serv, client)
    tCounterGet(t, stor, serv, client)
    tGaugeGet(t, stor, serv, client)

    
}
func tCounter(t *testing.T, stor *storageMock, serv *httptest.Server, client *http.Client){
     
    type want struct {
        code int
        res string
        val storage.Gauge
    }

    tests := []struct{
        URL string
        want want
    }{
        {
            URL: "",
            want: want{
                code:http.StatusNotFound,
                res: "404 page not found\n",
            },
        },
        {
            URL: "/update",
            want: want{
                code: http.StatusBadRequest,
                res: "Bad Request\n",
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
            URL: "/update/random/name/123",
            want: want{
                code:http.StatusBadRequest,
                res: "Bad Request\n",
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
        t.Run("testCounter " + serv.URL+test.URL,func(t *testing.T) {
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
            assert.Equal(t, test.want.val, stor.val)
        })
    }
}
func tGauge(t *testing.T, stor *storageMock, serv *httptest.Server, client *http.Client) {
     
    type want struct {
        code int
        res string
        val storage.Gauge
    }

    tests := []struct{
        URL string
        want want
    }{
        {
            URL: "",
            want: want{
                code:http.StatusNotFound,
                res: "404 page not found\n",
            },
        },
        {
            URL: "/update",
            want: want{
                code: http.StatusBadRequest, 
                res: "Bad Request\n",
            },
        },
        {
            URL: "/update/random/name/123",
            want: want{
                code: http.StatusBadRequest,
                res: "Bad Request\n",
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
            assert.Equal(t, test.want.val, stor.val)
        })
    }
}
func tCounterGet(t *testing.T, stor *storageMock, serv *httptest.Server, client *http.Client) {

    type want struct {
        code int
        res string
    }

    tests := []struct{
        URL string
        value storage.Counter
        name string
        want want
    }{
        {
            URL: "",
            want: want{
                code:http.StatusNotFound,
                res: "404 page not found\n",
            },
        },
        {
            URL: "/value/",
            want: want{
                code: http.StatusBadRequest,
                res: "Bad Request\n",
            },
        },
        {
            URL: "/value/random/123",
            want: want{
                code: http.StatusBadRequest, 
                res: "Bad Request\n",
            },
        },
        {
            URL: "/value/counter/",
            want: want{
                code:http.StatusNotFound,
                res: "Not Found\n",
            },
        },
        {
            URL: "/value/counter/name234",
            name: "name",
            want: want{
                code:http.StatusNotFound,
                res: "Unknown metric\n",
            },
        },
        {
            URL: "/value/counter/name",
            value: 233,
            name: "name",
            want: want{
                code:http.StatusOK,
                res: "233\n",
            },
        },
        {
            URL: "/value/counter/name/dfs",
            want: want{
                code:http.StatusBadRequest,
                res: "Bad Request\n",
            },
        },
    }
    for _, test := range tests {
        t.Run(test.URL, func(t *testing.T) {

            stor.val = storage.Gauge(test.value)

            stor.name = test.name

            res, err := client.Get(serv.URL+test.URL)

            require.NoError(t,err)

            assert.Equal(t, test.want.code, res.StatusCode)
            r, _ := io.ReadAll(res.Body)
            res.Body.Close()
            assert.Equal(t, test.want.res, string(r))
        })
    }
}
func tGaugeGet(t *testing.T, stor *storageMock, serv *httptest.Server, client *http.Client) {
    type want struct {
        code int
        res string
    }

    tests := []struct{
        URL string
        value storage.Gauge
        name string
        want want
    }{
        {
            URL: "",
            want: want{
                code:http.StatusNotFound,
                res: "404 page not found\n",
            },
        },
        {
            URL: "/value/",
            want: want{
                code: http.StatusBadRequest,
                res: "Bad Request\n",
            },
        },
        {
            URL: "/value/random/123",
            want: want{
                code: http.StatusBadRequest, 
                res: "Bad Request\n",
            },
        },
        {
            URL: "/value/gauge/",
            want: want{
                code:http.StatusNotFound,
                res: "Not Found\n",
            },
        },
        {
            URL: "/value/gauge/name234",
            name: "name",
            want: want{
                code:http.StatusNotFound,
                res: "Unknown metric\n",
            },
        },
        {
            URL: "/value/gauge/name",
            value: 233.213,
            name: "name",
            want: want{
                code:http.StatusOK,
                res: "233.213\n",
            },
        },
        {
            URL: "/value/gauge/name1",
            value: 0,
            name: "name1",
            want: want{
                code:http.StatusOK,
                res: "0\n",
            },
        },
        {
            URL: "/value/gauge/name2",
            value: -0,
            name: "name2",
            want: want{
                code:http.StatusOK,
                res: "0\n",
            },
        },
        {
            URL: "/value/gauge/name/dfs",
            want: want{
                code:http.StatusBadRequest,
                res: "Bad Request\n",
            },
        },
    }
    for _, test := range tests {
        t.Run(test.URL, func(t *testing.T) {

            stor.val = storage.Gauge(test.value)

            stor.name = test.name

            res, err := client.Get(serv.URL+test.URL)

            require.NoError(t,err)

            assert.Equal(t, test.want.code, res.StatusCode)
            r, _ := io.ReadAll(res.Body)
            res.Body.Close()
            assert.Equal(t, test.want.res, string(r))
        })
    }
}
