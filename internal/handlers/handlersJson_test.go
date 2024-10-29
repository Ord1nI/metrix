package handlers

import (
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	 gomock "github.com/golang/mock/gomock"

	// "go.uber.org/zap"

	"bytes"
	"encoding/json"
	"errors"

	// "io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ord1nI/metrix/internal/repo/metrics"

	"github.com/Ord1nI/metrix/mocks/mock_repo"
)


func ptrToInt(d int64) *int64 {
	return &d
}
func ptrToFloat(d float64) *float64 {
	return &d
}

func Test(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mR := mock_repo.NewMockRepo(ctrl)

	r := chi.NewRouter()
	r.Method(http.MethodPost, "/update/", APIFunc(UpdateJSON(mR)))
	r.Method(http.MethodPost, "/value/", APIFunc(GetJSON(mR)))

	TUpdateJSON(t, r, mR)
	TGetJSON(t, r, mR)

}

func TUpdateJSON(t *testing.T, r chi.Router, mr *mock_repo.MockRepo) {

	type want struct {
		responseM metrics.Metric
		response  string
		code      int
	}
	tests := []struct {
		metric metrics.Metric
		name   string
		want   want
	}{
		{
			name: "Test badReq",
			metric: metrics.Metric{
				ID:    "name",
				MType: "",
				Delta: ptrToInt(213),
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "Error while updating\n",
			},
		},
		{
			name: "Test badReq2",
			metric: metrics.Metric{
				ID:    "",
				MType: "counter",
				Delta: ptrToInt(213),
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "Error while updating\n",
			},
		},
		{
			name: "test gauge",
			metric: metrics.Metric{
				ID:    "gauge",
				MType: "gauge",
				Value: ptrToFloat(213),
			},
			want: want{
				code: http.StatusOK,
				responseM: metrics.Metric{
					ID:    "gauge",
					MType: "gauge",
					Value: ptrToFloat(213),
				},
			},
		},
		{
			name: "test counter",
			metric: metrics.Metric{
				ID:    "counter",
				MType: "counter",
				Delta: ptrToInt(1),
			},
			want: want{
				code: http.StatusOK,
				responseM: metrics.Metric{
					ID:    "counter",
					MType: "counter",
					Delta: ptrToInt(1),
				},
			},
		},
		{
			name: "test counter2",
			metric: metrics.Metric{
				ID:    "counter",
				MType: "counter",
				Delta: ptrToInt(1),
			},
			want: want{
				code: http.StatusOK,
				responseM: metrics.Metric{
					ID:    "counter",
					MType: "counter",
					Delta: ptrToInt(2),
				},
			},
		},
		{
			name: "test gauge2",
			metric: metrics.Metric{
				ID:    "gauge",
				MType: "gauge",
				Value: ptrToFloat(213),
			},
			want: want{
				code: http.StatusOK,
				responseM: metrics.Metric{
					ID:    "gauge",
					MType: "gauge",
					Value: ptrToFloat(213),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			buf := bytes.NewBuffer(nil)
			//
			err := json.NewEncoder(buf).Encode(&test.metric)

			if test.want.code == 200 {
				mr.EXPECT().Add(gomock.Any(),gomock.Any()).Return(nil)
				mr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil)
			} else {
				mr.EXPECT().Add(gomock.Any(),gomock.Any()).Return(errors.New("err"))
			}

			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/update/", buf)
			req.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			res := w.Result()

			if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
				var j metrics.Metric

				err = json.NewDecoder(res.Body).Decode(&j)

				require.NoError(t, err)

				assert.Equal(t, test.want.code, res.StatusCode)
			} else {
				assert.Equal(t, test.want.code, res.StatusCode)
				// b, err := io.ReadAll(res.Body)
				// require.NoError(t,err)
				// assert.Equal(t, test.want.response, string(b))
			}

			res.Body.Close()
		})
	}
}

func TGetJSON(t *testing.T, r chi.Router, mr *mock_repo.MockRepo) {

	type want struct {
		responseM metrics.Metric
		response  string
		code      int
	}
	tests := []struct {
		metric metrics.Metric
		name   string
		want   want
	}{
		{
			name: "Test badReq",
			metric: metrics.Metric{
				ID:    "some name",
				MType: "gauge",
			},
			want: want{
				code:     http.StatusNotFound,
				response: "Error while getting\n",
			},
		},
		{
			name: "test gauge",
			metric: metrics.Metric{
				ID:    "gauge",
				MType: "gauge",
			},
			want: want{
				code: http.StatusOK,
				responseM: metrics.Metric{
					ID:    "gauge",
					MType: "gauge",
					Value: ptrToFloat(213),
				},
			},
		},
		{
			name: "test counter",
			metric: metrics.Metric{
				ID:    "counter",
				MType: "counter",
			},
			want: want{
				code: http.StatusOK,
				responseM: metrics.Metric{
					ID:    "counter",
					MType: "counter",
					Delta: ptrToInt(2),
				},
			},
		},
		{
			name: "test gauge2",
			metric: metrics.Metric{
				ID:    "gauge",
				MType: "gauge",
			},
			want: want{
				code: http.StatusOK,
				responseM: metrics.Metric{

					MType: "gauge",
					Value: ptrToFloat(213),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if test.want.code == 200 {
				mr.EXPECT().Get(test.metric.ID, gomock.Any()).Return(nil)
			} else {
				mr.EXPECT().Get(test.metric.ID, gomock.Any()).Return(errors.New("err"))
			}

			buf := bytes.NewBuffer(nil)
			//
			err := json.NewEncoder(buf).Encode(&test.metric)

			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/value/", buf)
			req.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			res := w.Result()

			if res.Header.Get("Content-Type") == "application/json" {
				var j metrics.Metric

				err = json.NewDecoder(res.Body).Decode(&j)

				require.NoError(t, err)

				assert.Equal(t, test.want.code, res.StatusCode)
			} else {
				assert.Equal(t, test.want.code, res.StatusCode)
				// b, err := io.ReadAll(res.Body)
				// require.NoError(t,err)
				// assert.Equal(t, test.want.response, string(b))
			}

			res.Body.Close()
		})
	}
}
