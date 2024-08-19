package handlers 

import (
    "net/http"
    "errors"

    "github.com/Ord1nI/metrix/internal/repo"
    "github.com/Ord1nI/metrix/internal/repo/database"
)

func PingDB(r repo.Repo) APIFunc{
    hf := func(res http.ResponseWriter, req *http.Request) error {
        v, ok := r.(*database.Database)
        if ok {
            err := v.DB.PingContext(req.Context())

            if err != nil {
                return NewHandlerError(err, http.StatusInternalServerError)
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

