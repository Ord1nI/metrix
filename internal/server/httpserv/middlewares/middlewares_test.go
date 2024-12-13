package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/hex"
	"crypto/sha256"
	"crypto/hmac"

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

func TestSign(t *testing.T) {

	key := []byte("key")

	signedBody := []byte("uwu")

	signer := hmac.New(sha256.New, key)

	signer.Write(signedBody)

	Hash := hex.EncodeToString(signer.Sum(nil))

	type want struct {
		code int
	}

	tests := []struct {
		body []byte
		contentType string
		want        want
	}{
		// {
		// 	contentType: "",
		// 	want: want{
		// 		code: http.StatusOK,
		// 	},
		// },
		{
			contentType: "key",
			want: want {
				code: http.StatusNotFound,
			},
		},
		{
			contentType: "1234567890123456789012345678901234567890123456789012345678901234",
			want: want {
				code: http.StatusBadRequest,
			},
		},
		{
			contentType: Hash,
			want: want {
				code: http.StatusOK,
			},
		},
	}

	for v, test := range tests {
		t.Run(fmt.Sprintf("test %d", v), func(*testing.T) {

		recorder := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(signedBody))

		req.Header.Add("HashSHA256", test.contentType)

		f := SignMW(zap.NewNop().Sugar(), key)(http.HandlerFunc(handlerMock))

		f.ServeHTTP(recorder,req)

		assert.Equal(t, test.want.code, recorder.Code)

		})
	}
}

func headMock(res http.ResponseWriter, req *http.Request) {
	io.ReadAll(req.Body)
	req.Body.Close()
	io.ReadAll(req.Body)

	res.WriteHeader(http.StatusOK)
}


func TestHead(t *testing.T) {

	recorder := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("data")))

	f := HeadMW(zap.NewNop().Sugar())(http.HandlerFunc(headMock))

	f.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
