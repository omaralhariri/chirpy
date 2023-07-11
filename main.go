package main

import (
    "log"
    "strconv"
    "net/http"
)

type HealthCheck struct {}

type apiConfig struct {
    fileServerHits int
}

func main() {
    const port = "8080"

    mux := http.NewServeMux()
    logMux := middlewareLog(mux)
    corsMux := middlewareCors(logMux)
    
    apiCfg := &apiConfig {
        fileServerHits: 0,
    }

    mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
    mux.Handle("/app/assets/logo.png", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
    mux.Handle("/healthz", HealthCheck{})
    mux.Handle("/healthz2", http.HandlerFunc(health2))
    mux.Handle("/metrics", http.HandlerFunc(apiCfg.metricHandler))

    srv := &http.Server{
        Addr: ":" + port,
        Handler: corsMux,
    }

    log.Printf("Serving on port: %s\n", port)
    log.Fatal(srv.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "*")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileServerHits++

        next.ServeHTTP(w, r)
    })
}

func middlewareLog(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func (h HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("Ok"))
    return
}

func health2(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK 2"))
    return
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hits: " + strconv.Itoa(cfg.fileServerHits)))
    return
}
