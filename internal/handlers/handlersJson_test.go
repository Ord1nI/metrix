package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ord1nI/metrix/internal/myjson"
)

func (s *storageMock) AddMetric(m myjson.Metric) error{
    if m.ID == "" {
        return errors.New("error")
    }
    if m.MType == "counter" {
        if s.mtype == "counter" {
            s.val += float64(*m.Delta)
            s.name = m.ID
            return nil
        } else {
            s.val = float64(*m.Delta)
            s.name = m.ID
            s.mtype = "counter"
            return nil
        }
    }
    if m.MType == "gauge" {
        s.val = float64(*m.Value)
        s.name = m.ID
        s.mtype = "gauge"
        return nil
    }
    return errors.New("error")
}

func (s *storageMock) GetMetric(name, t  string) (*myjson.Metric, bool) {
    if name == s.name {
        switch t{
        case "counter":
            v := int64(s.val)
            return &myjson.Metric{
                ID: name,
                MType: t,
                Delta: &v,
            },true
        case "gauge":
            return &myjson.Metric{
                ID: name,
                MType: t,
                Value: &s.val,
            },true

        }
    }
    return nil, false
}

func ptrToInt(d int64) *int64 {
    return &d
}
func ptrToFloat(d float64) *float64 {
    return &d
}

func Test(t *testing.T) {
    r := chi.NewRouter()
    r.Method(http.MethodPost, "/update/", UpdateJSON(&storageMock{}))
    r.Method(http.MethodPost, "/value/",GetJSON(&storageMock{}))

    TUpdateJSON(t, r)

}

func TUpdateJSON(t *testing.T, r chi.Router) {

    type want struct {
        code int
        response string
        responseM myjson.Metric
    }
    tests := []struct{
        name string
        metric myjson.Metric
        want want
    }{
        {
            name: "Test badReq",
            metric: myjson.Metric{
                ID: "name",
                MType: "",
                Delta:ptrToInt(213),
            },
            want: want{
                code: http.StatusBadRequest,
                response: "error",
            },
        },
        {
            name: "Test badReq2",
            metric: myjson.Metric{
                ID: "",
                MType: "counter",
                Delta:ptrToInt(213),
            },
            want: want{
                code: http.StatusBadRequest,
                response: "error",
            },
        },
        {
            name: "test gauge",
            metric: myjson.Metric{
                ID: "gauge",
                MType: "gauge",
                Value:ptrToFloat(213),
            },
            want: want{
                code: http.StatusOK,
                responseM: myjson.Metric{
                    ID: "gauge",
                    MType: "gauge",
                    Value:ptrToFloat(213),
                },
            },
        },
        {
            name: "test counter",
            metric: myjson.Metric{
                ID: "counter",
                MType: "counter",
                Delta:ptrToInt(1),
            },
            want: want{
                code: http.StatusOK,
                responseM: myjson.Metric{
                    ID: "counter",
                    MType: "counter",
                    Delta:ptrToInt(1),
                },
            },
        },
        {
            name: "test counter2",
            metric: myjson.Metric{
                ID: "counter",
                MType: "counter",
                Delta:ptrToInt(1),
            },
            want: want{
                code: http.StatusOK,
                responseM: myjson.Metric{
                    ID: "counter",
                    MType: "counter",
                    Delta:ptrToInt(2),
                },
            },
        },
        {
            name: "test gauge2",
            metric: myjson.Metric{
                ID: "gauge",
                MType: "gauge",
                Value:ptrToFloat(213),
            },
            want: want{
                code: http.StatusOK,
                responseM: myjson.Metric{
                    ID: "gauge",
                    MType: "gauge",
                    Value:ptrToFloat(213),
                },
            },
        },
    }


    for _, test := range tests{
        t.Run(test.name, func(t *testing.T) {

            buf := bytes.NewBuffer(nil)
            //
            err := json.NewEncoder(buf).Encode(&test.metric)

            require.NoError(t,err)

            req := httptest.NewRequest(http.MethodPost, "/update/", buf)
            req.Header.Add("Content-Type", "application/json" )

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()

            if res.Header.Get("Content-Type") == "application/json" {
                var j myjson.Metric

                err = json.NewDecoder(res.Body).Decode(&j)

                require.NoError(t,err)

                assert.Equal(t, test.want.code, res.StatusCode)
                assert.Equal(t, test.want.responseM, j)
            } else {
                assert.Equal(t, test.want.code,res.StatusCode)
                b, err := io.ReadAll(res.Body)
                require.NoError(t,err)
                assert.Equal(t, test.want.response, string(b))
            }

            res.Body.Close()
        })
    }
}

func TGetJSON(t *testing.T, r chi.Router) {

    type want struct {
        code int
        response string
        responseM myjson.Metric
    }
    tests := []struct{
        name string
        metric myjson.Metric
        want want
    }{
        {
            name: "Test badReq",
            metric: myjson.Metric{
                ID:"some name",
                MType: "gauge",
            },
            want: want{
                code: http.StatusBadRequest,
                response: "error",
            },
        },
        {
            name: "test gauge",
            metric: myjson.Metric{
                ID:"gauge",
                MType:"gauge",
            },
            want: want{
                code: http.StatusOK,
                responseM: myjson.Metric{
                    ID: "gauge",
                    MType: "gauge",
                    Value:ptrToFloat(213),
                },
            },
        },
        {
            name: "test counter",
            metric:myjson.Metric{
                ID:"counter",
                MType: "counter",
            },
            want: want{
                code: http.StatusOK,
                responseM: myjson.Metric{
                    ID: "counter",
                    MType: "counter",
                    Delta:ptrToInt(2),
                },
            },
        },
        {
            name: "test gauge2",
            metric: myjson.Metric{
                ID:"gauge",
                MType:"gauge",
            },
            want: want{
                code: http.StatusOK,
                responseM: myjson.Metric{
                    ID: "gauge",
                    MType: "gauge",
                    Value:ptrToFloat(213),
                },
            },
        },
    }
    for _, test := range tests{
        t.Run(test.name, func(t *testing.T) {

            buf := bytes.NewBuffer(nil)
            //
            err := json.NewEncoder(buf).Encode(&test.metric)

            require.NoError(t,err)

            req := httptest.NewRequest(http.MethodPost, "/value/", buf)
            req.Header.Add("Content-Type", "application/json" )

            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            res := w.Result()

            if res.Header.Get("Content-Type") == "application/json" {
                var j myjson.Metric

                err = json.NewDecoder(res.Body).Decode(&j)

                require.NoError(t,err)

                assert.Equal(t, test.want.code, res.StatusCode)
                assert.Equal(t, test.want.responseM, j)
            } else {
                assert.Equal(t, test.want.code,res.StatusCode)
                b, err := io.ReadAll(res.Body)
                require.NoError(t,err)
                assert.Equal(t, test.want.response, string(b))
            }

            res.Body.Close()
        })
    }
}
