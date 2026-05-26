package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/pz4-monitoring/internal/metrics"
	"example.com/pz4-monitoring/internal/student"
)

type Handler struct {
	repo *student.Repo
}

func NewHandler(repo *student.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawID := strings.TrimPrefix(r.URL.Path, "/students/")
	if rawID == "" || rawID == r.URL.Path {
		http.Error(w, "student id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		http.Error(w, "invalid student id", http.StatusBadRequest)
		return
	}

	// Доп. задание 1: счётчик запросов по student_id
	metrics.StudentRequestsTotal.WithLabelValues(rawID).Inc()

	// Доп. задание 2: histogram только для /students/{id}
	start := time.Now()
	defer func() {
		metrics.StudentRequestDuration.Observe(time.Since(start).Seconds())
	}()

	st, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}
