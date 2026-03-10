package models

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type CreateTaskRequest struct {
	Title string `json:"title" example:"Go to the gym"`
}

type PatchTaskRequest struct {
	Done bool `json:"done" example:"true"`
}
