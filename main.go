package main

import (
    "log"
    "net/http"
)

func main () {
    const filepathRoot = "."
    const port = "8080"

    mux := http.NewServeMux()
    mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
    mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
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
