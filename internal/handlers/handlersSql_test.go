package handlers

import (
	"net/http"
	"net/http/httptest"
	"bytes"
	"testing"

	"github.com/Ord1nI/metrix/internal/repo"

	"github.com/stretchr/testify/assert"
)

func TestPingDB(t *testing.T) {
	var db repo.Repo
	f := PingDB(db)


	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet,"/", bytes.NewBuffer(nil))


	f.ServeHTTP(recorder,req)

	assert.Equal(t,http.StatusBadRequest, recorder.Result().StatusCode)
}
