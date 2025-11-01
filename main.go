package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/jifpbj/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening sql: %v", err)
	}
	dbQueries := database.New(db)
	user, err := dbQueries.CreateUser(context.Background(), "a@gmail.com")
	fmt.Printf("user: %v", user.Email)

	const filepathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	handler := http.StripPrefix(
		"/app",
		apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", handler)
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
	})

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/users", handlerCreateUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
