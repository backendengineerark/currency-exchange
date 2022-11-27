package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	db, err := PrepareDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetExchange(db))
	})

	mux.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetAllExchangesSaved(db))
	})

	http.ListenAndServe(":8080", mux)
}
