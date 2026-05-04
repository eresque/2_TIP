package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"example.com/order-service/internal/order"
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
	// Variant 4: read user-service address from env
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8081"
	}

	repo := order.NewRepo()
	client := order.NewUserServiceClient(userServiceURL)
	handler := order.NewHandler(repo, client)

	mux := http.NewServeMux()

	mux.HandleFunc("/orders/by-user/", handler.GetOrdersByUser)
	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/full") {
			handler.GetOrderWithUser(w, r)
			return
		}
		handler.GetOrderByID(w, r)
	})

	addr := ":8082"
	log.Println("order-service started on", addr)
	log.Println("user-service URL:", userServiceURL)
	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
