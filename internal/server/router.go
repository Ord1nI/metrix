package server

import (
	"github.com/go-chi/chi/v5"

	"net/http"
    "time"
    
	"github.com/Ord1nI/metrix/internal/handlers"
	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo"
)

func updateGaugeRoute(sugar logger.Logger, stor repo.Adder, BackoffSchedule []time.Duration) func(r chi.Router){
    return func(r chi.Router) {
        // ANY /update/gauge/
        r.HandleFunc("/", handlers.NotFound) 
        // POST /update/gauge/name/123
        r.Method(http.MethodPost, "/{name}/{val}", handlers.Make(sugar, handlers.UpdateGauge(stor), BackoffSchedule))
        // ANY /update/gauge/name/123/adsf
        r.HandleFunc("/{name}/{val}/*", handlers.BadRequest)    
    }
}

func updateCounterRoute(sugar logger.Logger, stor repo.Adder, BackoffSchedule []time.Duration) func(r chi.Router){
    return func(r chi.Router) {
        // ANY /update/gauge/
        r.HandleFunc("/", handlers.NotFound)                    
        // POST /update/gauge/name/123
        r.Method(http.MethodPost, "/{name}/{val}", 
            handlers.Make(sugar, handlers.UpdateCounter(stor), BackoffSchedule))
        // ANY /update/gauge/name/123/adsf
        r.HandleFunc("/{name}/{val}/*", handlers.BadRequest)  
    }
}
func valueGaugeRoute(sugar logger.Logger, stor repo.Getter, BackoffSchedule []time.Duration) func(r chi.Router){
    return func(r chi.Router) {
        //ANY /value/gauge/
        r.HandleFunc("/", handlers.NotFound)  
        //GET /value/gauge/name
        r.Method(http.MethodGet, "/{name}",
            handlers.Make(sugar, handlers.GetGauge(stor), BackoffSchedule))
        //ANY /value/gauge/name/asa
        r.HandleFunc("/{name}/*", handlers.BadRequest)      
    }
}
func valueCounterRoute(sugar logger.Logger, stor repo.Getter, BackoffSchedule []time.Duration) func(r chi.Router){
    return func(r chi.Router) {
        //ANY /value/counter/
        r.HandleFunc("/", handlers.NotFound)                   
        //GET /value/counter/name
        r.Method(http.MethodGet,"/{name}", 
            handlers.Make(sugar, handlers.GetCounter(stor), BackoffSchedule))
        //ANY /value/counter/name/qew
        r.HandleFunc("/{name}/*", handlers.BadRequest)      
    }
}

func (s *Server) InitRouter(middlewares ...func(http.Handler)http.Handler) {
    s.Router = CreateRouter(s.Logger, s.Repo, s.Config.BackoffSchedule, middlewares...)
}

func CreateRouter(log logger.Logger, re repo.Repo, BackoffSchedule []time.Duration, middlewares ...func(http.Handler)http.Handler) chi.Router{
    r := chi.NewRouter()

    for _, i := range middlewares {
        r.Use(i)
    }


    // GET /
    r.Method(http.MethodGet, "/", 
        handlers.Make(log, handlers.MainPage(re), BackoffSchedule))

    r.Method(http.MethodGet, "/ping", 
        handlers.Make(log, handlers.PingDB(re), BackoffSchedule))

    r.Method(http.MethodPost, "/updates/", 
        handlers.Make(log, handlers.UpdatesJSON(re), BackoffSchedule))

    r.Route("/update", func(r chi.Router) {
        // POST /pudate/
        r.Method(http.MethodPost, "/", 
            handlers.Make(log, handlers.UpdateJSON(re), BackoffSchedule))
        // ANY /update/*
        r.HandleFunc("/*", handlers.BadRequest)                      
        // ANY /update/gauge/*
        r.Route("/gauge", updateGaugeRoute(log, re, BackoffSchedule))
        // Any /update/counter/*
        r.Route("/counter", updateCounterRoute(log, re, BackoffSchedule))     
        
    })

    r.Route("/value", func(r chi.Router) {
        r.Method(http.MethodPost, "/", 
            handlers.Make(log, handlers.GetJSON(re), BackoffSchedule))
        // Any /value/
        r.HandleFunc("/*", handlers.BadRequest)
        // ANY /value/gauge/*
        r.Route("/gauge", valueGaugeRoute(log, re, BackoffSchedule))        
        // ANY /value/counter/*
        r.Route("/counter", valueCounterRoute(log, re, BackoffSchedule))   
    })
    return r
}
