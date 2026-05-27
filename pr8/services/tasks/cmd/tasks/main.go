package main

import (
	"log"
	"net/http"
	"os"

	tasks "example.com/tasks/internal"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", tasks.HealthHandler)

	addr := ":" + port
	log.Println("tasks service started on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
