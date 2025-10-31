package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	apiCfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hits reset back to 0"))
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
	w.Header().Add("Content-Type", "text/html")
	_, err := w.Write([]byte(htmlResponse))
	if err != nil {
		fmt.Printf("error writing request : %v", err)
	}
}

func (apiCfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
