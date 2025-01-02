package main

import (
    "log"
    "net/http"
)

func main () {
    mux := http.NewServeMux()
    srv := &http.Server {
        Addr: ":8080",
        Handler: mux,
    }

    log.Println("Starting server on :8080...")
    log.Fatal(srv.ListenAndServe())

}
