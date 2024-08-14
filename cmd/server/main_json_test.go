package main

import (

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"io"
	"net/http"
    "bytes"
	"net/http/httptest"
    "encoding/json"
	"testing"
    "reflect"

	"github.com/Ord1nI/metrix/internal/myjson"
	"github.com/Ord1nI/metrix/internal/storage"
)

func Test(t *testing.T) {
    stor := storage.NewMemStorage()

    r := CreateRouter(stor)

    serv := httptest.NewServer(r)
    client := serv.Client()

    tUpdateJSON(t,serv, client)
    tGetJSON(t, serv, client)
}

func ptrFloat(d float64) *float64{
    return &d
}
func ptrInt(d int64) *int64{
    return &d
}

func tUpdateJSON(t *testing.T, serv *httptest.Server, client *http.Client) {
    type want struct {
        code int
        res string
        metric myjson.Metric
    }

    tests := []struct{
        name string
        metric myjson.Metric
        want want
    }{
        {
            name: "name1",
            metric: myjson.Metric{
                ID: "name",
                MType: "gauge",
                Value: ptrFloat(51.34),
            },
            want: want{
                metric: myjson.Metric{
                    ID: "name",
                    MType: "gauge",
                    Value: ptrFloat(51.34),
                },
                code:http.StatusOK,
                res: "",
            },
        },
        {
            name: "with out metric name",
            metric: myjson.Metric{
                ID: "",
                MType: "gauge",
                Value: ptrFloat(51.34),
            },
            want: want{
                code:http.StatusBadRequest,
                res: "metric must hame name",
            },
        },
        {
            name: "with out metric type",
            metric: myjson.Metric{
                ID: "name",
                MType: "",
                Value: ptrFloat(51.34),
            },
            want: want{
                code: http.StatusBadRequest,
                res: "bad type",
            },
        },
        {
            name: "name3",
            metric: myjson.Metric{
                ID: "name",
                MType: "counter",
                Delta: ptrInt(123),
            },
            want: want{
                metric: myjson.Metric{
                    ID: "name",
                    MType: "counter",
                    Delta: ptrInt(123),
                },
                code:http.StatusOK,
                res: "",
            },
        },
        {
            name: "name3",
            metric: myjson.Metric{
                ID: "name",
                MType: "counter",
                Delta: ptrInt(123),
            },
            want: want{
                code:http.StatusOK,
                metric: myjson.Metric{
                    ID: "name",
                    MType: "counter",
                    Delta: ptrInt(246),
                },
                res: "",
            },
        },
    }

    for _, test := range tests {
        t.Run("testJson " + test.name, func(t *testing.T) {

            data, err := json.Marshal(test.metric)
            require.NoError(t,err)
            buf := bytes.NewBuffer(data)

            res, err := client.Post(serv.URL + "/update/", "application/json", buf)

            require.NoError(t,err)

            assert.Equal(t, test.want.code, res.StatusCode)

            if res.Header.Get("Content-Type") == "application/json" {
                var metric myjson.Metric

                data, _ = io.ReadAll(res.Body)
                json.Unmarshal(data, &metric)

                assert.Equal(t, test.want.metric, metric)
            } else {
                r, _ := io.ReadAll(res.Body)
                assert.Equal(t, test.want.res, string(r))
            }
            res.Body.Close()
        })
    }

}
func tGetJSON(t *testing.T, serv *httptest.Server, client *http.Client) {
    type want struct {
        code int
        res string
        metric myjson.Metric
    }

    tests := []struct{
        name string
        metric myjson.Metric
        want want
    }{
        {
            name: "name1",
            metric: myjson.Metric{
                ID: "name",
                MType: "gauge",
            },
            want: want{
                code:http.StatusOK,
                res: "",
                metric: myjson.Metric{
                    ID: "name",
                    MType: "gauge",
                    Value: ptrFloat(51.34),
                },
            },
        },
        {
            name: "without metric name",
            metric: myjson.Metric{
                ID: "",
                MType: "gauge",
            },
            want: want{
                code:http.StatusNotFound,
                res: "Cant find this metric",
            },
        },
        {
            name: "without metric type",
            metric: myjson.Metric{
                ID: "name",
                MType: "",
            },
            want: want{
                code:http.StatusNotFound,
                res: "Cant find this metric",
            },
        },
        {
            name: "name2",
            metric: myjson.Metric{
                ID: "name",
                MType: "counter",
            },
            want: want{
                metric: myjson.Metric{
                    ID: "name",
                    MType: "counter",
                    Delta: ptrInt(246),
                },
                code:http.StatusOK,
                res: "",
            },
        },
    }

    for _, test := range tests {
        t.Run("testJson " + test.name, func(t *testing.T) {

            data, err := json.Marshal(test.metric)
            require.NoError(t,err)
            buf := bytes.NewBuffer(data)

            res, err := client.Post(serv.URL + "/value/", "application/json", buf)

            require.NoError(t,err)

            assert.Equal(t, test.want.code, res.StatusCode)

            if res.Header.Get("Content-Type") == "application/json" {
                var metric myjson.Metric

                data, _ = io.ReadAll(res.Body)
                json.Unmarshal(data, &metric)

                assert.True(t, reflect.DeepEqual(test.want.metric, metric))
            } else {
                r, _ := io.ReadAll(res.Body)
                assert.Equal(t, test.want.res, string(r))
            }
            res.Body.Close()
        })
    }
}
