package config

import "os"

type Config struct {
	Addr     string
	CertFile string
	KeyFile  string
	DSN      string
}

func New() Config {
	return Config{
		Addr:     getEnv("ADDR", ":8443"),
		CertFile: getEnv("CERT_FILE", "certs/server.crt"),
		KeyFile:  getEnv("KEY_FILE", "certs/server.key"),
		DSN:      getEnv("DSN", "postgres://postgres:postgres@localhost:5432/study_security?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
