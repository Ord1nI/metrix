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
        r.Handle("/", logger.HandlerLoggingFn(handlers.NotFound, sugar)) // ANY /update/gauge/
        r.Method(http.MethodPost, "/{name}/{val}", logger.HandlerLogging(handlers.UpdateGauge(stor), sugar))     // POST /update/gauge/name/123
        r.Handle("/{name}/{val}/*", logger.HandlerLoggingFn(handlers.BadRequest, sugar))    // ANY /update/gauge/name/123/adsf
    }
}

func updateCounterRoute(stor storage.Adder) func(r chi.Router){
    return func(r chi.Router) {
        r.Handle("/", logger.HandlerLoggingFn(handlers.NotFound, sugar))                    // ANY /update/gauge/
        r.Method(http.MethodPost, "/{name}/{val}", logger.HandlerLogging(handlers.UpdateCounter(stor), sugar))     // POST /update/gauge/name/123
        r.Handle("/{name}/{val}/*", logger.HandlerLoggingFn(handlers.BadRequest, sugar))  // ANY /update/gauge/name/123/adsf
    }
}
func valueGaugeRoute(stor storage.Getter) func(r chi.Router){
    return func(r chi.Router) {
        r.Handle("/", logger.HandlerLoggingFn(handlers.NotFound, sugar))  //ANY /value/gauge/
        r.Method(http.MethodGet, "/{name}", logger.HandlerLogging(handlers.GetGauge(stor), sugar))           //GET /value/gauge/name
        r.Handle("/{name}/*", logger.HandlerLoggingFn(handlers.BadRequest, sugar))      //ANY /value/gauge/name/asa
    }
}
func valueCounterRoute(stor storage.Getter) func(r chi.Router){
    return func(r chi.Router) {
        r.Handle("/", logger.HandlerLoggingFn(handlers.NotFound, sugar))                   //ANY /value/counter/
        r.Method(http.MethodGet,"/{name}", logger.HandlerLogging(handlers.GetCounter(stor), sugar))         //GET /value/counter/name
        r.Handle("/{name}/*", logger.HandlerLoggingFn(handlers.BadRequest,sugar))      //ANY /value/counter/name/qew
    }
}

func CreateRouter(stor *storage.MemStorage) *chi.Mux{

    r := chi.NewRouter()

    r.Method(http.MethodGet, "/", logger.HandlerLogging(handlers.GetAllMetrics(stor), sugar))                  //POST localhost:/

    r.Route("/update", func(r chi.Router) {
        r.Handle("/*", logger.HandlerLoggingFn(handlers.BadRequest, sugar))                      // ANY /update/


        r.Route("/gauge", updateGaugeRoute(stor))         // ANY /update/gauge/*

        r.Route("/counter", updateCounterRoute(stor))     // Any /update/counter/*
        
    })

    r.Route("/value", func(r chi.Router) {
        r.Handle("/*", logger.HandlerLoggingFn(handlers.BadRequest,sugar))            // Any /value/

        r.Route("/gauge", valueGaugeRoute(stor))        

        r.Route("/counter", valueCounterRoute(stor))   
    })

    return r
}
