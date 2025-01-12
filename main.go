package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Todo struct {
	Id        string `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

var (
	todos     []Todo
	todoMutex sync.Mutex
)

func generateRandomId() string {
	return uuid.New().String()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:  "ok",
		Message: "Api is working",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)

	case "POST":
		var newTodo Todo
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// fmt.Println(string(body))
		err = json.Unmarshal(body, &newTodo)
		if err != nil || newTodo.Task == "" {
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}

		fmt.Printf("Unmarshalled Todo: %+v\n", newTodo)

		newTodo.Id = generateRandomId()
		todoMutex.Lock()
		todos = append(todos, newTodo)
		todoMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo)

		// fmt.Printf("Todos: %+v\n", todos)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func todoByIdYHandler(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

func main() {

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/todos", todosHandler)
	http.HandleFunc("/todos/", todoByIdYHandler)

	fmt.Println("App is running in port 3000")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting app", err)
	}
}
