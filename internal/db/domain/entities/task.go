package entities

import "github.com/google/uuid"

type Task struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Done    bool      `json:"done"`
	UserID  uuid.UUID `json:"user_id"`
}

func NewTask(title, content string, userID uuid.UUID) *Task {
	return &Task{
		ID:      uuid.New(),
		Title:   title,
		Content: content,
		Done:    false,
		UserID:  userID,
	}
}
