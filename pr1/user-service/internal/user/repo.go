package user

import "errors"

var ErrUserNotFound = errors.New("user not found")

type Repo struct {
	data map[int64]User
}

func NewRepo() *Repo {
	return &Repo{
		data: map[int64]User{
			1: {ID: 1, Name: "Иван Иванов", Email: "ivan@example.com"},
			2: {ID: 2, Name: "Мария Петрова", Email: "maria@example.com"},
			3: {ID: 3, Name: "Алексей Сидоров", Email: "alex@example.com"},
		},
	}
}

func (r *Repo) GetByID(id int64) (User, error) {
	u, ok := r.data[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

// Variant 2: return all users
func (r *Repo) GetAll() []User {
	users := make([]User, 0, len(r.data))
	for _, u := range r.data {
		users = append(users, u)
	}
	return users
}
