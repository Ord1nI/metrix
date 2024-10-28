package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ord1nI/metrix/internal/repo"
	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
)

func TestPingDB(t *testing.T) {
	var db repo.Repo

	f := PingDB(db)
	r := chi.NewRouter()

	r.Method(http.MethodGet, "/", f)

	req := httptest.NewRequest(http.MethodGet,"/", bytes.NewBuffer(nil))

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder,req)

	res := recorder.Result()

	defer res.Body.Close()

	assert.Equal(t,http.StatusBadRequest, recorder.Result().StatusCode)

}
