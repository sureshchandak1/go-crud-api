package storage

import "github.com/sureshchandak1/go-crud-api/internal/types"

type Storage interface {
	CreateStudent(name, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
}
