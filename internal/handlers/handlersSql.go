package handlers 

import (
    "net/http"

    "github.com/Ord1nI/metrix/internal/repo"
    "github.com/Ord1nI/metrix/internal/repo/database"
    "github.com/Ord1nI/metrix/internal/logger"
)

func PingDB(l logger.Logger, r repo.Repo) http.Handler{
    hf := func(res http.ResponseWriter, req *http.Request) {
        v, ok := r.(*database.Database)
        if ok {
            err := v.Db.PingContext(req.Context())

            if err != nil {
                l.Errorln(err)
                http.Error(res, err.Error(), http.StatusInternalServerError)
                return
            } else {
                res.WriteHeader(http.StatusOK)
                res.Write([]byte("successfuly connect to database"))
                return
            }
        } 
        http.Error(res, "no database in use", http.StatusInternalServerError)
    }
    return http.HandlerFunc(hf)
}

