package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func handlerMock(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func TestCompressor(t *testing.T) {
	type CT struct {
		contentType string
		val         string
	}
	type want struct {
		contentType CT
	}

	tests := []struct {
		want        want
		contentType []CT
		code        int
	}{
		{
			contentType: []CT{CT{"Content-Encoding", "gzip"}},
			code:        500,
		},
		{
			contentType: []CT{CT{"Accept-Encoding", "gzip"}, CT{"Accept", "html"}},
			want: want{
				contentType: CT{"Content-Encoding", "gzip"},
			},
			code: 200,
		},
		{
			contentType: []CT{CT{"", ""}},
			code:        200,
		},
	}

	for v, test := range tests {
		t.Run(fmt.Sprintf("test %d", v), func(*testing.T) {

			buf := bytes.NewBuffer(nil)
			req, _ := http.NewRequest(http.MethodGet, "/", buf)
			recorder := httptest.NewRecorder()

			for _, v := range test.contentType {
				req.Header.Add(v.contentType, v.val)
			}

			handler := CompressorMW(zap.NewNop().Sugar())(http.HandlerFunc(handlerMock))
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, test.code, recorder.Code)
			assert.Equal(t, test.want.contentType.val, recorder.Header().Get(test.want.contentType.contentType))
		})
	}
}
