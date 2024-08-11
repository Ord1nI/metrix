package main

import (
    "github.com/go-chi/chi/v5"

    "net/http"
    "github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/handlers"
    "github.com/Ord1nI/metrix/internal/logger"
)


func updateGaugeRoute(stor storage.Adder) func(r chi.Router){
    return func(r chi.Router) {
        r.HandleFunc("/", handlers.NotFound) // ANY /update/gauge/
        r.Method(http.MethodPost, "/{name}/{val}", handlers.UpdateGauge(stor))     // POST /update/gauge/name/123
        r.HandleFunc("/{name}/{val}/*", handlers.BadRequest)    // ANY /update/gauge/name/123/adsf
    }
}

func updateCounterRoute(stor storage.Adder) func(r chi.Router){
    return func(r chi.Router) {
        r.HandleFunc("/", handlers.NotFound)                    // ANY /update/gauge/
        r.Method(http.MethodPost, "/{name}/{val}", handlers.UpdateCounter(stor))     // POST /update/gauge/name/123
        r.HandleFunc("/{name}/{val}/*", handlers.BadRequest)  // ANY /update/gauge/name/123/adsf
    }
}
func valueGaugeRoute(stor storage.Getter) func(r chi.Router){
    return func(r chi.Router) {
        r.HandleFunc("/", handlers.NotFound)  //ANY /value/gauge/
        r.Method(http.MethodGet, "/{name}", handlers.GetGauge(stor))           //GET /value/gauge/name
        r.HandleFunc("/{name}/*", handlers.BadRequest)      //ANY /value/gauge/name/asa
    }
}
func valueCounterRoute(stor storage.Getter) func(r chi.Router){
    return func(r chi.Router) {
        r.HandleFunc("/", handlers.NotFound)                   //ANY /value/counter/
        r.Method(http.MethodGet,"/{name}", handlers.GetCounter(stor))         //GET /value/counter/name
        r.HandleFunc("/{name}/*", handlers.BadRequest)      //ANY /value/counter/name/qew
    }
}

func CreateRouter(stor *storage.MemStorage) *chi.Mux{

    r := chi.NewRouter()

    r.Use(logger.HandlerLogging(sugar))

    r.Method(http.MethodGet, "/", handlers.MainPage(stor))                  //POST localhost:/

    r.Route("/update", func(r chi.Router) {
        r.HandleFunc("/*", handlers.BadRequest)                      // ANY /update/


        r.Route("/gauge", updateGaugeRoute(stor))         // ANY /update/gauge/*

        r.Route("/counter", updateCounterRoute(stor))     // Any /update/counter/*
        
    })

    r.Route("/value", func(r chi.Router) {
        r.HandleFunc("/*", handlers.BadRequest)            // Any /value/

        r.Route("/gauge", valueGaugeRoute(stor))        

        r.Route("/counter", valueCounterRoute(stor))   
    })

    return r
}
