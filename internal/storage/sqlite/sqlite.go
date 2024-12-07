package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sureshchandak1/go-crud-api/internal/config"
	"github.com/sureshchandak1/go-crud-api/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil

}

func (s *Sqlite) GetStudents() ([]types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return []types.Student{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []types.Student{}, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return []types.Student{}, err
		}

		students = append(students, student)
	}

	return students, nil

}

func (s *Sqlite) UpdateStudentById(id int64, name, email string, age int) error {

	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite) DeleteStudentById(id int64) (types.Student, error) {
	return types.Student{}, nil
}
