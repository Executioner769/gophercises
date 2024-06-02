package db

type Store interface {
	CreateTask(*Task) (*Task, error)
	GetTasks() ([]*Task, error)
	CompleteTask(int) (*Task, error)
	DeleteTask(int) (*Task, error)
}
