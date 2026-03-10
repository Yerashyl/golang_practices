package handlers

import (
	"assignment1/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

var tasks = map[int]models.Task{
	1: {ID: 1, Title: "Go touch grass", Done: false},
	2: {ID: 2, Title: "Don't be cooked", Done: false},
}

var nextID = 3

// GetTasks godoc
// @Summary      Get tasks
// @Description  Get a list of all tasks or a specific task by ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   query      int  false  "Task ID"
// @Success      200  {array}    models.Task
// @Failure      400  {object}   map[string]string
// @Failure      404  {object}   map[string]string
// @Security     ApiKeyAuth
// @Router       /tasks [get]
func GetTasks(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	if idStr == "" {

		list := []models.Task{}

		for _, t := range tasks {
			list = append(list, t)
		}

		writeJSON(w, list, http.StatusOK)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, ok := tasks[id]
	if !ok {
		writeError(w, "task not found", http.StatusNotFound)
		return
	}

	writeJSON(w, task, http.StatusOK)
}

// CreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task with a title
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body models.CreateTaskRequest true "Task title"
// @Success      201  {object}   models.Task
// @Failure      400  {object}   map[string]string
// @Security     ApiKeyAuth
// @Router       /tasks [post]
func CreateTask(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Title string `json:"title"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		writeError(w, "title must not be empty", http.StatusBadRequest)
		return
	}

	task := models.Task{
		ID:    nextID,
		Title: req.Title,
		Done:  false,
	}

	tasks[nextID] = task
	nextID++

	writeJSON(w, task, http.StatusCreated)
}

// PatchTask godoc
// @Summary      Update a task's status
// @Description  Update the "done" status of an existing task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   query      int  true  "Task ID"
// @Param        request body models.PatchTaskRequest true "Task status"
// @Success      200  {object}   models.Task
// @Failure      400  {object}   map[string]string
// @Failure      404  {object}   map[string]string
// @Security     ApiKeyAuth
// @Router       /tasks [patch]
func PatchTask(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, ok := tasks[id]
	if !ok {
		writeError(w, "task not found", http.StatusNotFound)
		return
	}

	var req struct {
		Done bool `json:"done"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, "invalid json", http.StatusBadRequest)
		return
	}

	task.Done = req.Done
	tasks[id] = task

	writeJSON(w, task, http.StatusOK)
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a task by ID
// @Tags         tasks
// @Produce      json
// @Param        id   query      int  true  "Task ID"
// @Success      204  "No Content"
// @Failure      400  {object}   map[string]string
// @Failure      404  {object}   map[string]string
// @Security     ApiKeyAuth
// @Router       /tasks [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, "invalid id", http.StatusBadRequest)
		return
	}

	_, ok := tasks[id]
	if !ok {
		writeError(w, "task not found", http.StatusNotFound)
		return
	}

	delete(tasks, id)

	w.WriteHeader(http.StatusNoContent)
}

// GetExternalTodos godoc
// @Summary      Get external todos
// @Description  Fetch todos from an external JSON placeholder API
// @Tags         external
// @Produce      json
// @Success      200  {array}    interface{}
// @Failure      500  {object}   map[string]string
// @Router       /external/todos [get]
func GetExternalTodos(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")
	if err != nil {
		writeError(w, "external api error", 500)
		return
	}

	defer resp.Body.Close()

	var data interface{}

	json.NewDecoder(resp.Body).Decode(&data)

	writeJSON(w, data, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, code int) {

	writeJSON(w, map[string]string{
		"error": msg,
	}, code)
}
