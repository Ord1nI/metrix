package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/repo/storage"
)

func Test(t *testing.T) {

    logger := zap.NewNop()
    sugar = logger.Sugar()
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
        metric metrics.Metric
    }

    tests := []struct{
        name string
        metric metrics.Metric
        want want
    }{
        {
            name: "name1",
            metric: metrics.Metric{
                ID: "name",
                MType: "gauge",
                Value: ptrFloat(51.34),
            },
            want: want{
                metric: metrics.Metric{
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
            metric: metrics.Metric{
                ID: "",
                MType: "gauge",
                Value: ptrFloat(51.34),
            },
            want: want{
                code:http.StatusBadRequest,
                res: "Error while updating\n",
            },
        },
        {
            name: "with out metric type",
            metric: metrics.Metric{
                ID: "name",
                MType: "",
                Value: ptrFloat(51.34),
            },
            want: want{
                code: http.StatusBadRequest,
                res: "Error while updating\n",
            },
        },
        {
            name: "name3",
            metric: metrics.Metric{
                ID: "name",
                MType: "counter",
                Delta: ptrInt(123),
            },
            want: want{
                metric: metrics.Metric{
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
            metric: metrics.Metric{
                ID: "name",
                MType: "counter",
                Delta: ptrInt(123),
            },
            want: want{
                code:http.StatusOK,
                metric: metrics.Metric{
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
                var metric metrics.Metric

                data, _ = io.ReadAll(res.Body)
                json.Unmarshal(data, &metric)

                assert.Equal(t, test.want.metric, metric)
            } else {
                // r, _ := io.ReadAll(res.Body)
                // assert.Equal(t, test.want.res, string(r))
            }
            res.Body.Close()
        })
    }

}
func tGetJSON(t *testing.T, serv *httptest.Server, client *http.Client) {
    type want struct {
        code int
        res string
        metric metrics.Metric
    }

    tests := []struct{
        name string
        metric metrics.Metric
        want want
    }{
        {
            name: "name1",
            metric: metrics.Metric{
                ID: "name",
                MType: "gauge",
            },
            want: want{
                code:http.StatusOK,
                res: "",
                metric: metrics.Metric{
                    ID: "name",
                    MType: "gauge",
                    Value: ptrFloat(51.34),
                },
            },
        },
        {
            name: "without metric name",
            metric: metrics.Metric{
                ID: "",
                MType: "gauge",
            },
            want: want{
                code:http.StatusNotFound,
                res: "Metric not found\n",
            },
        },
        {
            name: "without metric type",
            metric: metrics.Metric{
                ID: "name",
                MType: "",
            },
            want: want{
                code:http.StatusNotFound,
                res: "Metric not found\n",
            },
        },
        {
            name: "name2",
            metric: metrics.Metric{
                ID: "name",
                MType: "counter",
            },
            want: want{
                metric: metrics.Metric{
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
                var metric metrics.Metric

                data, _ = io.ReadAll(res.Body)
                json.Unmarshal(data, &metric)

                assert.True(t, reflect.DeepEqual(test.want.metric, metric))
            } else {
                // r, _ := io.ReadAll(res.Body)
                // assert.Equal(t, test.want.res, string(r))
            }
            res.Body.Close()
        })
    }
}
