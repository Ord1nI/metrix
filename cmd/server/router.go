package main

import (
	"github.com/go-chi/chi/v5"

	"net/http"
    
	"github.com/Ord1nI/metrix/internal/handlers"
	"github.com/Ord1nI/metrix/internal/repo"
)


func updateGaugeRoute(stor repo.Adder) func(r chi.Router){
    return func(r chi.Router) {
        // ANY /update/gauge/
        r.HandleFunc("/", handlers.NotFound) 
        // POST /update/gauge/name/123
        r.Method(http.MethodPost, "/{name}/{val}", handlers.UpdateGauge(sugar,stor))             
        // ANY /update/gauge/name/123/adsf
        r.HandleFunc("/{name}/{val}/*", handlers.BadRequest)    
    }
}

func updateCounterRoute(stor repo.Adder) func(r chi.Router){
    return func(r chi.Router) {
        // ANY /update/gauge/
        r.HandleFunc("/", handlers.NotFound)                    
        // POST /update/gauge/name/123
        r.Method(http.MethodPost, "/{name}/{val}", handlers.UpdateCounter(sugar,stor))     
        // ANY /update/gauge/name/123/adsf
        r.HandleFunc("/{name}/{val}/*", handlers.BadRequest)  
    }
}
func valueGaugeRoute(stor repo.Getter) func(r chi.Router){
    return func(r chi.Router) {
        //ANY /value/gauge/
        r.HandleFunc("/", handlers.NotFound)  
        //GET /value/gauge/name
        r.Method(http.MethodGet, "/{name}", handlers.GetGauge(sugar, stor))           
        //ANY /value/gauge/name/asa
        r.HandleFunc("/{name}/*", handlers.BadRequest)      
    }
}
func valueCounterRoute(stor repo.Getter) func(r chi.Router){
    return func(r chi.Router) {
        //ANY /value/counter/
        r.HandleFunc("/", handlers.NotFound)                   
        //GET /value/counter/name
        r.Method(http.MethodGet,"/{name}", handlers.GetCounter(sugar, stor))         
        //ANY /value/counter/name/qew
        r.HandleFunc("/{name}/*", handlers.BadRequest)      
    }
}

func CreateRouter(stor repo.Repo, middlewares ...func(http.Handler)http.Handler) *chi.Mux{

    r := chi.NewRouter()

    for _, i := range middlewares {
        r.Use(i)
    }


    // GET /
    r.Method(http.MethodGet, "/", handlers.MainPage(sugar, stor))                  

    r.Method(http.MethodGet, "/ping", handlers.PingDB(sugar, stor))                  

    r.Method(http.MethodPost, "/updates/", handlers.UpdatesJSON(sugar, stor))

    r.Route("/update", func(r chi.Router) {
        r.Method(http.MethodPost, "/", handlers.UpdateJSON(sugar, stor))
        // ANY /update/
        r.HandleFunc("/*", handlers.BadRequest)                      
        // ANY /update/gauge/*
        r.Route("/gauge", updateGaugeRoute(stor))
        // Any /update/counter/*
        r.Route("/counter", updateCounterRoute(stor))     
        
    })

    r.Route("/value", func(r chi.Router) {
        r.Method(http.MethodPost, "/", handlers.GetJSON(sugar, stor))
        // Any /value/
        r.HandleFunc("/*", handlers.BadRequest)
        // ANY /value/gauge/*
        r.Route("/gauge", valueGaugeRoute(stor))        
        // ANY /value/counter/*
        r.Route("/counter", valueCounterRoute(stor))   
    })

    return r
}
