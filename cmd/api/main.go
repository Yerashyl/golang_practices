package main

import (
	"assignment1/internal/handlers"
	"assignment1/internal/middleware"
	"log"
	"net/http"

	_ "assignment1/docs" // Ignore unused import error for docs since swagger requires it

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           Tasks API
// @version         1.0
// @description     This is a sample server for assignment 1.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY

func main() {

	mux := http.NewServeMux()

	taskRouter := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			handlers.GetTasks(w, r)

		case http.MethodPost:
			handlers.CreateTask(w, r)

		case http.MethodPatch:
			handlers.PatchTask(w, r)

		case http.MethodDelete:
			handlers.DeleteTask(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/tasks",
		middleware.Logging(
			middleware.RequestID(
				middleware.APIKey(taskRouter),
			),
		),
	)

	mux.HandleFunc("/external/todos", handlers.GetExternalTodos)

	// Swagger endpoint
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	log.Println("server started :8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}
