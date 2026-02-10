package storage

type Storage interface {
	createStudent(name string, email string, age int) (int64, error)
}
