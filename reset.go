package main

import (
    "fmt"
    "net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Store(0)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hits reset to 0")))
}
