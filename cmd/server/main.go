package main

import (
    "net/http"
    "strings"
    "strconv"
    "fmt"
)

type MemStorage struct{
    gauge map[string] float64
    counter map[string] int64
}

var data MemStorage = MemStorage{
    gauge: make(map[string]float64),
    counter: make(map[string]int64),
}

func updateGauge(res http.ResponseWriter, req *http.Request) {

    if req.Method != http.MethodGet {
        http.Error(res, "Only Get requests in allowd", http.StatusMethodNotAllowed)
        return
    }

    url := strings.Split(req.URL.Path, "/")[3:]

    if len(url) != 2 {
        http.Error(res, "Incorect", http.StatusBadRequest)
        return
    }

    name := url[0]
    val,err := strconv.ParseFloat(url[1], 64)

    if err != nil {
        http.Error(res, "Incorect value", http.StatusBadRequest)
        return
    }

    data.gauge[name] = val
    res.WriteHeader(http.StatusOK)
}

func updateCounter(res http.ResponseWriter, req *http.Request) {

    if req.Method != http.MethodGet {
        http.Error(res, "Only Get requests in allowd", http.StatusMethodNotAllowed)
        return
    }

    url := strings.Split(req.URL.Path, "/")[3:]

    if len(url) != 2 {
        http.Error(res, "Incorect", http.StatusBadRequest)
        return
    }
    
    name := url[0]
    val, err := strconv.ParseInt(url[1], 10, 64)

    if err != nil {
        http.Error(res, "Incorect value", http.StatusBadRequest)
        return
    }

    data.counter[name] += val
    res.WriteHeader(http.StatusOK)
    fmt.Println(data.counter)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusBadRequest)})
    mux.HandleFunc(`/update/gauge/`, updateGauge)
    mux.HandleFunc(`/update/counter/`, updateCounter)

    err := http.ListenAndServe(`:8080`, mux)
    if err != nil {
        panic(err)
    }
}
