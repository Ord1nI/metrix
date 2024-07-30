package main

import (
    "flag"
    "net/http"
    "github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/handlers"
)

var fIpStr = flag.String("a",":8080","enter IP format ip:port")


func main() {
    flag.Parse()


    stor := storage.NewEmptyStorage()

    r := CreateRouter(stor)

    r.Get("/", handlers.GetAllMetrics(stor))                  //POST localhost:/


    err := http.ListenAndServe(*fIpStr, r)
    if err != nil {
        panic(err)
    }
}
