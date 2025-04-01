package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
func handleAllTasks(res http.ResponseWriter, req *http.Request) {

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)

}

func handleTaskGet(res http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")

	t, ok := tasks[id]

	if ok {
		resp, err := json.Marshal(t)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(resp)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusBadRequest)

}

func handleTaskPost(res http.ResponseWriter, req *http.Request) {

	var t Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &t); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[t.ID] = t

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

}

func handleTaskDelete(res http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	_, ok := tasks[id]

	if ok {

		delete(tasks, id)
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusBadRequest)

}

func main() {

	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", handleAllTasks)
	r.Post("/tasks", handleTaskPost)
	r.Get("/tasks/{id}", handleTaskGet)
	r.Delete("/tasks/{id}", handleTaskDelete)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
