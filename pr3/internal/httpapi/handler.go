package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"example.com/pz3-logging/internal/student"
	"go.uber.org/zap"
)

type Handler struct {
	repo *student.Repo
	log  *zap.Logger
}

func NewHandler(repo *student.Repo, log *zap.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.log.Warn("method not allowed for health endpoint",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.log.Debug("health endpoint called")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.log.Warn("method not allowed for student endpoint",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/students/")
	if path == "" || path == r.URL.Path {
		h.log.Warn("student id is missing",
			zap.String("path", r.URL.Path),
		)
		http.Error(w, "student id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		h.log.Warn("invalid student id",
			zap.String("raw_id", path),
			zap.Error(err),
		)
		http.Error(w, "invalid student id", http.StatusBadRequest)
		return
	}

	// Доп. задание 3: debug-лог перед обращением к репозиторию
	h.log.Debug("looking up student in repository",
		zap.Int64("student_id", id),
	)

	st, err := h.repo.GetByID(id)
	if err != nil {
		h.log.Error("student not found",
			zap.Int64("student_id", id),
			zap.Error(err),
		)
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}

	h.log.Info("student returned successfully",
		zap.Int64("student_id", st.ID),
		zap.String("group", st.Group),
	)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}

// Доп. задание 4: POST /students — создание нового студента
func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.log.Warn("method not allowed for create student endpoint",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input student.Student
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("failed to decode request body",
			zap.Error(err),
		)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.FullName == "" || input.Group == "" || input.Email == "" {
		h.log.Warn("validation failed: missing required fields",
			zap.String("full_name", input.FullName),
			zap.String("group", input.Group),
			zap.String("email", input.Email),
		)
		http.Error(w, "full_name, group and email are required", http.StatusBadRequest)
		return
	}

	created := h.repo.Create(input)
	h.log.Info("student created successfully",
		zap.Int64("student_id", created.ID),
		zap.String("full_name", created.FullName),
		zap.String("group", created.Group),
	)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(created)
}
