package models

type Task struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Overdue     bool   `json:"overdue"`
	Completed   bool   `json:"completed"`
}
