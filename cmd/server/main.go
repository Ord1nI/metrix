package main

import (
    "net/http"
    "github.com/Ord1nI/metrix/cmd/handlers"
    "github.com/Ord1nI/metrix/cmd/storage"
)



func main() {

    stor := storage.NewEmptyStorage()

    mux := http.NewServeMux()
    mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusBadRequest)})
    mux.HandleFunc(`/update/gauge/`, handlers.UpdateGauge(stor))
    mux.HandleFunc(`/update/counter/`, handlers.UpdateCounter(stor))

    err := http.ListenAndServe(`:8080`, mux)
    if err != nil {
        panic(err)
    }
}
