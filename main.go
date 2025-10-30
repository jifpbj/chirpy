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

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{}

	handler := http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))))

	mux := http.NewServeMux()
	mux.Handle("/app/", handler)
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
	})

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	apiCfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Reset App OK!"))
	if err != nil {
		fmt.Printf("error writing statusText: %s", err)
	}
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		fmt.Printf("error writing statusText: %s", err)
	}
}

func (apiCfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	v := apiCfg.fileserverHits.Load()

	htmlResponse := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`, v)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(htmlResponse))
	if err != nil {
		fmt.Printf("error writing request : %v", err)
	}
}

func (apiCfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { apiCfg.fileserverHits.Add(1); next.ServeHTTP(w, r) })
}
