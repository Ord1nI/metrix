package handlers 

import (
    "database/sql"
    "net/http"
)

func PingDB(db *sql.DB) http.Handler{
    hf := func(res http.ResponseWriter, req *http.Request) {
        err := db.PingContext(req.Context())

        if err != nil {
            http.Error(res, err.Error(), http.StatusInternalServerError)
            return
        } else {
            res.WriteHeader(http.StatusOK)
            res.Write([]byte("successfuly connect to database"))
        }
    }
    return http.HandlerFunc(hf)
}

