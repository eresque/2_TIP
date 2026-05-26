package httpapi

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"example.com/pz5-security/internal/student"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Handler struct {
	repo *student.Repo
	stmt *sql.Stmt
}

func NewHandler(repo *student.Repo, stmt *sql.Stmt) *Handler {
	return &Handler{
		repo: repo,
		stmt: stmt,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"scheme": "https",
	})
}

func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawID := r.URL.Query().Get("id")
	if rawID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// Allow-list валидация: id должен быть положительным числом (доп. задание 4).
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id: must be a positive integer", http.StatusBadRequest)
		return
	}

	var st student.Student
	err = h.stmt.QueryRow(id).Scan(&st.ID, &st.FullName, &st.StudyGroup, &st.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}

// GetStudentByEmail — поиск студента по email (доп. задание 3).
// Allow-list валидация формата email реализована через regexp (доп. задание 4).
func (h *Handler) GetStudentByEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	// Allow-list валидация: email должен соответствовать стандартному формату.
	if !emailRegex.MatchString(email) {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	st, err := h.repo.GetByEmail(email)
	if err != nil {
		if err == student.ErrStudentNotFound {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}
