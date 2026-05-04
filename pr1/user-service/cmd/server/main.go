package main

import (
	"log"
	"net/http"
	"time"

	"example.com/user-service/internal/user"
)

// Variant 3: logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	repo := user.NewRepo()
	handler := user.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", handler.GetUsers)
	mux.HandleFunc("/users/", handler.GetUserByID)

	addr := ":8081"
	log.Println("user-service started on", addr)
	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
