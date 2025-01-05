package main

import (
    "fmt"
    "log"
    "net/http"
    "sync/atomic"
)

type apiConfig struct {
    fileserverHits atomic.Int32
}

func main () {
    const filepathRoot = "."
    const port = "8080"

    apiCfg := apiConfig{
        fileserverHits: atomic.Int32{},
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            http.Redirect(w, r, "/app/", http.StatusFound)
            return
        }
        http.NotFound(w, r)
    })

    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
    mux.HandleFunc("/healthz", handlerReadiness)
    mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
    mux.HandleFunc("/reset", apiCfg.handlerReset)
    
    srv := &http.Server {
        Addr: ":" + port,
        Handler: mux,
    }

    log.Printf("Serving files form %s on port: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())

}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Cache-Control") == "no-cache" {
            w.Header().Set("Cache-Control", "no-cache")
        }
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}
