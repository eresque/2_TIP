package student

import (
	"errors"
	"sync"

	"example.com/pz2-grpc/gen/studentpb"
)

var ErrStudentNotFound = errors.New("student not found")

type Repository struct {
	mu     sync.Mutex
	data   map[int64]*studentpb.Student
	nextID int64
}

func NewRepository() *Repository {
	return &Repository{
		nextID: 4,
		data: map[int64]*studentpb.Student{
			1: {
				Id:             1,
				FullName:       "Иванов Иван Иванович",
				Group:          "ИВБО-01-25",
				Email:          "ivanov@example.com",
				Specialization: "Информационная безопасность",
			},
			2: {
				Id:             2,
				FullName:       "Петрова Мария Сергеевна",
				Group:          "ИВБО-02-25",
				Email:          "petrova@example.com",
				Specialization: "Программная инженерия",
			},
			3: {
				Id:             3,
				FullName:       "Сидоров Алексей Андреевич",
				Group:          "ИВБО-03-25",
				Email:          "sidorov@example.com",
				Specialization: "Прикладная математика",
			},
		},
	}
}

func (r *Repository) GetByID(id int64) (*studentpb.Student, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	st, ok := r.data[id]
	if !ok {
		return nil, ErrStudentNotFound
	}
	return st, nil
}

func (r *Repository) GetAll() []*studentpb.Student {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := make([]*studentpb.Student, 0, len(r.data))
	for _, st := range r.data {
		list = append(list, st)
	}
	return list
}

func (r *Repository) Create(fullName, group, email, specialization string) *studentpb.Student {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextID
	r.nextID++
	st := &studentpb.Student{
		Id:             id,
		FullName:       fullName,
		Group:          group,
		Email:          email,
		Specialization: specialization,
	}
	r.data[id] = st
	return st
}
