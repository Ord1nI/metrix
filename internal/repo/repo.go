package repo

import (
	"encoding/json"
)

type Repo interface {
	json.Marshaler
	Add(name string, val interface{}) error
	Get(name string, val interface{}) error
	Close() error
}
