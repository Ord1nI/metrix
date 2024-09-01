package handlers

import (
	"errors"
	"net/http"

	"github.com/Ord1nI/metrix/internal/repo"
	"github.com/Ord1nI/metrix/internal/repo/database"
)

func PingDB(r repo.Repo) APIFunc {
	hf := func(res http.ResponseWriter, req *http.Request) error {
		v, ok := r.(*database.Database)
		if ok {
			err := v.Ping()

			if err != nil {
				return NewHandlerError(errors.Join(err, errors.New("ping error")), http.StatusInternalServerError)
			} else {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte("successfuly connect to database"))
				return nil
			}
		}
		return NewHandlerError(errors.New("doesnt use database"), http.StatusBadRequest)
	}
	return APIFunc(hf)
}
