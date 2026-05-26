package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"example.com/pz5-security/internal/config"
	"example.com/pz5-security/internal/httpapi"
	"example.com/pz5-security/internal/student"
)

func main() {
	cfg := config.New()

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := student.NewRepo(db)

	stmt, err := repo.PrepareGetByID()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	handler := httpapi.NewHandler(repo, stmt)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/students", handler.GetStudentByID)
	mux.HandleFunc("/students/by-email", handler.GetStudentByEmail)

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := &http.Server{
		Addr:      cfg.Addr,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	// Доп. задание 1: HTTP-сервер на 8080 перенаправляет все запросы на HTTPS.
	go func() {
		redirectMux := http.NewServeMux()
		redirectMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			target := "https://localhost:8443" + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})
		log.Println("HTTP redirect server started on http://localhost:8080")
		if err := http.ListenAndServe(":8080", redirectMux); err != nil {
			log.Println("HTTP redirect server error:", err)
		}
	}()

	log.Println("HTTPS server started on https://localhost:8443")
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}
