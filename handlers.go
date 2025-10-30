package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	type JSONinput struct {
		Body string `json:"body"`
	}

	type JSONoutput struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	input := JSONinput{}
	output := JSONoutput{true, ""}

	err := decoder.Decode(&input)
	if err != nil {
		output.Error = "Something went wrong"
		output.Valid = false
		w.WriteHeader(500)
	}

	if len(input.Body) > 140 {
		output.Error = "Chirp is too long"
		output.Valid = false
		w.WriteHeader(400)
	}

	dat, err := json.Marshal(output)
	if err != nil {
		fmt.Printf("error marshaling output json: %v", err)
		w.WriteHeader(500)
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(dat)
	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("error writing data: %v", err)
	}
}

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
