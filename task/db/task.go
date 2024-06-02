package db

import "time"

type Task struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

func NewTask(message string) *Task {
	return &Task{
		Message:   message,
		Done:      false,
		CreatedAt: time.Now(),
	}
}
