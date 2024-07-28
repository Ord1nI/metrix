package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type storageMock struct {
}

    func (storageMock) GetGauge(name string) (float64, error) {
        return 0, nil
    }
    func (storageMock) AddGauge(name string, val float64) {
    }
    func (storageMock) GetCounter(name string) (int64, error) {
        return 0, nil
    }
    func (storageMock) AddCounter(name string, val int64) {
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
            name: "Test 1",
            reqURL: "http://fuckintsite.com/update/gauge/",
            want: want{
                code: http.StatusBadRequest,
                response: "Bad request\n",
            },
        },
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
    for _, test := range tests{
        t.Run(test.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            UpdateGauge(storageMock{})(w, req)

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
            name: "Test 1",
            reqURL: "http://fuckintsite.com/update/counter/",
            want: want{
                code: http.StatusBadRequest,
                response: "Bad request\n",
            },
        },
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
    for _, test := range tests{
        t.Run(test.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodPost, test.reqURL, nil)

            w := httptest.NewRecorder()

            UpdateCounter(storageMock{})(w, req)

            res := w.Result()
            assert.Equal(t,test.want.code, res.StatusCode)

            resBody, err := io.ReadAll(res.Body)

            require.NoError(t,err)
            
            assert.Equal(t,test.want.response, string(resBody))
            res.Body.Close()
        })
    }
}
