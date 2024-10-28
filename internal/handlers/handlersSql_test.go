package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ord1nI/metrix/internal/repo"

	"github.com/stretchr/testify/assert"
)

func TestPingDB(t *testing.T) {
	var db repo.Repo

	f := PingDB(db)

	req := httptest.NewRequest(http.MethodGet,"/", nil)

	recorder := httptest.NewRecorder()

	f.ServeHTTP(recorder,req)

	res := recorder.Result()

	defer res.Body.Close()

	assert.Equal(t,http.StatusBadRequest, recorder.Result().StatusCode)

}
