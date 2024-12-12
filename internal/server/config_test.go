package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("ADDRESS", "10.5.12.1")
	os.Setenv("STORE_INTERVAL", "5")
	os.Setenv("FILE_STORAGE_PATH", "5")
	os.Setenv("RESTORE", "true")
	os.Setenv("DATABASE_DSN", "5")
	os.Setenv("KEY", "asdf")

	serv, _ := New()
	assert.Equal(t, "10.5.12.1", serv.Config.Address)
	assert.Equal(t, int64(5), serv.Config.StoreInterval)
	assert.Equal(t, "5", serv.Config.FileStoragePath)
	assert.Equal(t, true, serv.Config.Restore)
	assert.Equal(t, "5", serv.Config.DBdsn)
	assert.Equal(t, "asdf", serv.Config.Key)
}
