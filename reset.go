package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
        _, err := w.Write([]byte("Reset is only allowed in dev mode"))
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Failed to write response", err)
            return
        }
		return
	}

	cfg.fileserverHits.Store(0)
    err := cfg.db.Reset(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to reset database", err)
        return
    }
	w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hits reset to 0 and database reset to initial state")
}
