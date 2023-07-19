package main

import (
    "net/http"
    "fmt"
    "html/template"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileServerHits++

        next.ServeHTTP(w, r)
    })
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits)))
    return
}

func (cfg *apiConfig) metricHandlerHtml(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/html;")
    t, err := template.ParseFiles("admin/metrics.html")
    if err != nil {
        panic(err)
    }
    w.WriteHeader(http.StatusOK)
    t.Execute(w, cfg.fileServerHits)
    return
}
