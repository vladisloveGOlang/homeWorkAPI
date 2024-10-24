package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...
func getHandler(w http.ResponseWriter, r *http.Request) {

	bytes, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	var task Task
	err = json.Unmarshal(bytes, &task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var task Task

	err = json.Unmarshal(bytes, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if tasks[task.ID].ID == "" {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Задачи с таким номером не существует"))
		return
	}

	delete(tasks, task.ID)
	w.WriteHeader(http.StatusAccepted)

}

func putchHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var task Task

	err = json.Unmarshal(bytes, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.WriteHeader(http.StatusOK)

}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getHandler)
	r.Post("/tasks", postHandler)
	r.Delete("/tasks", deleteHandler)
	r.Patch("/tasks", putchHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
