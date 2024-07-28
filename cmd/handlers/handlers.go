package handlers

import (
    "net/http"
    "strconv"
    "strings"
    "github.com/Ord1nI/metrix/cmd/storage"
)


func UpdateGauge(s storage.Repositories) func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        if req.Method != http.MethodPost {
            http.Error(res, "Only POST requests in allowd", http.StatusMethodNotAllowed)
            return
        }

        url := strings.Split(req.URL.Path, "/")[3:]

        if len(url) != 2 {
            http.Error(res, "Bad request", http.StatusBadRequest)
            return
        }

        name := url[0]
        val,err := strconv.ParseFloat(url[1], 64)

        if err != nil {
            http.Error(res, "Incorect metric value", http.StatusBadRequest)
            return
        }

        s.AddGauge(name,val)
        res.WriteHeader(http.StatusOK)
    }
}

func UpdateCounter(s storage.Repositories) func(res http.ResponseWriter, req *http.Request){
return func(res http.ResponseWriter, req *http.Request) {
        if req.Method != http.MethodPost {
            http.Error(res, "Only POST requests in allowd", http.StatusMethodNotAllowed)
            return
        }

        url := strings.Split(req.URL.Path, "/")[3:]

        if len(url) != 2 {
            http.Error(res, "Bad request", http.StatusBadRequest)
            return
        }
        
        name := url[0]
        val, err := strconv.ParseInt(url[1], 10, 64)

        if err != nil {
            http.Error(res, "Incorect metric value", http.StatusBadRequest)
            return
        }

        s.AddCounter(name, val)
        res.WriteHeader(http.StatusOK)
    }
}
