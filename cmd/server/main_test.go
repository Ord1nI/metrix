package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ord1nI/metrix/internal/repo/storage"
	"github.com/Ord1nI/metrix/internal/repo/metrics"
)

func TestMain(t *testing.T) {
    stor := storage.NewMemStorage()

    r := CreateRouter(stor)

    serv := httptest.NewServer(r)
    client := serv.Client()

    tCounter(t, stor, serv, client)
    tGauge(t, stor, serv, client)
    tCounterGet(t, serv, client)
    tGaugeGet(t, serv, client)
}

func tCounter(t *testing.T, stor *storage.MemStorage, serv *httptest.Server, client *http.Client){
     
    type want struct {
        code int
        res string
        val metrics.Counter
    }

    tests := []struct{
        URL string
        name string
        want want
    }{
        {
            URL: "/update",
            want: want{
                code: http.StatusBadRequest,
                res: "Not json request\n",
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
                res: "Error while updating\n",
            },
        },
        {
            URL: "/update/counter/name/123.213",
            name: "name",
            want: want{
                code:http.StatusBadRequest,
                res: "Error while updating\n",
                val: 0,
            },
        },
        {
            URL: "/update/counter/name1/-123",
            name: "name1",
            want: want{
                code:http.StatusOK,
                res: "",
                val: -123,
            },
        },
        {
            URL: "/update/counter/name1/123",
            name: "name1",
            want: want{
                code:http.StatusOK,
                res: "",
                val: 0,
            },
        },
    }
    for _, test := range tests {
        t.Run("testCounter " + serv.URL+test.URL,func(t *testing.T) {

            res, err := client.Post(serv.URL+test.URL,"text/plain",nil)

            require.NoError(t,err)
            defer res.Body.Close()

            assert.Equal(t, test.want.code, res.StatusCode)

            // r, _ := io.ReadAll(res.Body)
            // res.Body.Close()
            // assert.Equal(t, test.want.res, string(r))
            val, ok := stor.Counter.Get(test.name)
            if ok {
                assert.Equal(t, test.want.val, val)
            }
        })
    }
}
func tGauge(t *testing.T, stor *storage.MemStorage, serv *httptest.Server, client *http.Client) {
     
    type want struct {
        code int
        res string
        val metrics.Gauge
    }

    tests := []struct{
        URL string
        name string
        want want
    }{
        {
            URL: "/update",
            want: want{
                code: http.StatusBadRequest, 
                res: "Not json request\n",
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
                res: "Error while updating\n",
            },
        },
        {
            URL: "/update/gauge/name/123.213",
            name: "name",
            want: want{
                code:http.StatusOK,
                res: "",
                val: 123.213,
            },
        },
        {
            URL: "/update/gauge/name1/-123.213",
            name: "name1",
            want: want{
                code:http.StatusOK,
                res: "",
                val: -123.213,
            },
        },
    }
    for _, test := range tests {
        t.Run(serv.URL+test.URL,func(t *testing.T) {
            res, err := client.Post(serv.URL+test.URL,"text/plain",nil)

            require.NoError(t,err)
            defer res.Body.Close()
            assert.Equal(t, test.want.code, res.StatusCode)

            // r, _ := io.ReadAll(res.Body)
            // res.Body.Close()
            // assert.Equal(t, test.want.res, string(r))
            v, ok := stor.Gauge.Get(test.name)
            if ok {
                assert.Equal(t, test.want.val, v)
            }
        })
    }
}
func tCounterGet(t *testing.T,  serv *httptest.Server, client *http.Client) {

    type want struct {
        code int
        res string
    }

    tests := []struct{
        URL string
        name string
        want want
    }{
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
            name: "name234",
            want: want{
                code:http.StatusNotFound,
                res: "Metric not found\n",
            },
        },
        {
            URL: "/value/counter/name1",
            name: "name",
            want: want{
                code:http.StatusOK,
                res: "0\n",
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

            res, err := client.Get(serv.URL+test.URL)

            require.NoError(t,err)
            defer res.Body.Close()

            assert.Equal(t, test.want.code, res.StatusCode)
            // r, _ := io.ReadAll(res.Body)
            // res.Body.Close()
            // assert.Equal(t, test.want.res, string(r))
        })
    }
}
func tGaugeGet(t *testing.T,  serv *httptest.Server, client *http.Client) {
    type want struct {
        code int
        res string
    }

    tests := []struct{
        URL string
        name string
        want want
    }{
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
                res: "Metric not found\n",
            },
        },
        {
            URL: "/value/gauge/name",
            name: "name",
            want: want{
                code:http.StatusOK,
                res: "123.213\n",
            },
        },
        {
            URL: "/value/gauge/name1",
            name: "name1",
            want: want{
                code:http.StatusOK,
                res: "-123.213\n",
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

            res, err := client.Get(serv.URL+test.URL)

            require.NoError(t,err)
            defer res.Body.Close()

            assert.Equal(t, test.want.code, res.StatusCode)
            // r, _ := io.ReadAll(res.Body)
            // res.Body.Close()
            // assert.Equal(t, test.want.res, string(r))
        })
    }
}
