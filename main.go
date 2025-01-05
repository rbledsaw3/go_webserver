package main

import (
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
    fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
    rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            http.Redirect(w, r, "/app/", http.StatusFound)
            return
        }
        http.NotFound(w, r)
    })

    mux.HandleFunc("GET /", rootHandler)
    mux.Handle("GET /app/", fsHandler) 
    mux.HandleFunc("GET /healthz", handlerReadiness)
    mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
    mux.HandleFunc("POST /reset", apiCfg.handlerReset)

    mux.HandleFunc("POST /", handerMethodNotAllowed)
    mux.HandleFunc("POST /app/", handerMethodNotAllowed)
    mux.HandleFunc("POST /healthz", handerMethodNotAllowed)
    mux.HandleFunc("POST /metrics", handerMethodNotAllowed)
    mux.HandleFunc("GET /reset", handerMethodNotAllowed)
    
    srv := &http.Server {
        Addr: ":" + port,
        Handler: mux,
    }

    log.Printf("Serving files form %s on port: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())

}

func handerMethodNotAllowed (w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusMethodNotAllowed)
    w.Write([]byte("Method not allowed"))
}
