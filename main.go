package main

import (
    "fmt"
    "log"
    "net/http"
)

type apiConfig struct {
    fileserverHits atomic.Int32
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

func main () {
    const filepathRoot = "."
    const port = "8080"
    apiCfg := &apiConfig{}

    mux := http.NewServeMux()
    mux.Handle("/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))))
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "Hits: %d\n", apiCfg.fileserverHits.Load())
    })
    mux.HandleFunc("/reset", func(w http.ResponseWriter, req *http.Request) {
        apiCfg.fileserverHits.Store(0)
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    srv := &http.Server {
        Addr: ":" + port,
        Handler: mux,
    }

    log.Printf("Serving files form %s on port: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())

}
