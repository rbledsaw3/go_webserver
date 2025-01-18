package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "sync/atomic"

    "github.com/rbledsaw3/go_webserver/internal/database"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

type apiConfig struct {
    fileserverHits atomic.Int32
    db             *database.Queries
    platform       string
    jwtSecret      string
}

func main() {
    const filepathRoot = "."
    const port = "8080"

    godotenv.Load()
    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL must be set")
    }
    platform := os.Getenv("PLATFORM")
    if platform == "" {
        log.Fatal("PLATFORM must be set")
    }
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET env var must be set")
    }

    dbConn, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Error opening database: %s", err)
    }
    dbQueries := database.New(dbConn)

    apiCfg := apiConfig{
        fileserverHits: atomic.Int32{},
        db:             dbQueries,
        platform:       platform,
        jwtSecret:      jwtSecret,
    }

    mux := http.NewServeMux()
    fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
    mux.Handle("/app/", fsHandler)

    mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
    mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
    mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
    mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
    mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
    mux.HandleFunc("GET /api/healthz", handlerReadiness)
    mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
    mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
    mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
    mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
    mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: mux,
    }

    log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())
}
