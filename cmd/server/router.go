package main

import (
    "github.com/go-chi/chi/v5"
    "github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/handlers"
)

func updateGaugeRoute(stor storage.Adder) func(r chi.Router){
    return func(r chi.Router) {
        r.HandleFunc("/", handlers.NotFound)                    // ANY /update/gauge/
        r.Post("/{name}/{val}", handlers.UpdateGauge(stor))     // POST /update/gauge/name/123
        r.HandleFunc("/{name}/{val}/*", handlers.BadRequest)    // ANY /update/gauge/name/123/adsf
    }
}

func updateCounterRoute(stor storage.Adder) func(r chi.Router){
    return func(r chi.Router) {
        r.HandleFunc("/", handlers.NotFound)                    // ANY /update/gauge/
        r.Post("/{name}/{val}", handlers.UpdateCounter(stor))     // POST /update/gauge/name/123
        r.HandleFunc("/{name}/{val}/*", handlers.BadRequest)  // ANY /update/gauge/name/123/adsf
    }
}
func valueGaugeRoute(stor storage.Getter) func(r chi.Router){
    return func(r chi.Router) {
        r.HandleFunc("/", handlers.NotFound)                //ANY /value/gauge/
        r.Get("/{name}", handlers.GetGauge(stor))           //GET /value/gauge/name
        r.HandleFunc("/{name}/*", handlers.BadRequest)      //ANY /value/gauge/name/asa
    }
}
func valueCounterRoute(stor storage.Getter) func(r chi.Router){
    return func(r chi.Router) {
        r.HandleFunc("/", handlers.NotFound)                //ANY /value/counter/
        r.Get("/{name}", handlers.GetCounter(stor))         //GET /value/counter/name
        r.HandleFunc("/{name}/*", handlers.BadRequest)      //ANY /value/counter/name/qew
    }
}

func CreateRouter(stor storage.GetterAdder) *chi.Mux{

    r := chi.NewRouter()


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
