//Package repo contains interface to collect metrics.
package repo

import (
	"encoding/json"
)

//Repo base storage interface.
type Repo interface {
	json.Marshaler
    //Add function to add metric.
	Add(name string, val interface{}) error
    //Get function to get metric.
	Get(name string, val interface{}) error
    //Close function exist to close database connection not used in map.
	Close() error
}
