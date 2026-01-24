package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})

		r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message":"Hello from Chi","status":"ok"}`))
		})

		r.Post("/json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// Simple echo for benchmark purposes
			w.Write([]byte(`{"echo":"ok"}`))
		})

		r.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id":"` + id + `","name":"User ` + id + `"}`))
		})
	})

	log.Println("Chi server starting on :3004")
	log.Fatal(http.ListenAndServe(":3004", r))
}
