package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Струтура задачи
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// Список задач
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

/*
Обработчик для получения всех задач

Обработчик должен вернуть все задачи, которые хранятся в мапе.

- Конечная точка /tasks.

- Метод GET.

- При успешном запросе сервер должен вернуть статус 200 OK.

- При ошибке сервер должен вернуть статус 500 Internal Server Error.

Во всех обработчиках тип контента Content-Type — application/json.
*/
func getTasks(w http.ResponseWriter, r *http.Request) {
	// Сериализация данных
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Тип контента и статус
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Запись данных в тело ответа
	w.Write(resp)
}

/*
Обработчик для отправки задачи на сервер

Обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.

- Конечная точка /tasks.

- Метод POST.

- При успешном запросе сервер должен вернуть статус 201 Created.

- При ошибке сервер должен вернуть статус 400 Bad Request.

Во всех обработчиках тип контента Content-Type — application/json.
*/
func postTask(w http.ResponseWriter, r *http.Request) {
	// Переменная для новой задачи
	var task Task

	// Чтение POST-запроса
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Запись тела запроса в переменную
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавление задачи в список
	tasks[task.ID] = task

	// Тип контента и статус
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

/*
Обработчик должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе.

В мапе ключами являются ID задач. Вспомните, как проверить, есть ли ключ в мапе. Если такого ID нет, верните соответствующий статус.

- Конечная точка /tasks/{id}.

- Метод GET.

- При успешном выполнении запроса сервер должен вернуть статус 200 OK.

- В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.

Во всех обработчиках тип контента Content-Type — application/json.
*/
func getTaskById(w http.ResponseWriter, r *http.Request) {
	// Значение параметра из URL
	id := chi.URLParam(r, "id")

	// Поиск задачи по значению параметра
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// Сериализация данных
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Тип контента и статус
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Запись данных в тело ответа
	w.Write(resp)
}

/*
Обработчик удаления задачи по ID

Обработчик должен удалить задачу из мапы по её ID. Здесь так же нужно сначала проверить, есть ли задача с таким ID в мапе, если нет вернуть соответствующий статус.

- Конечная точка /tasks/{id}.

- Метод DELETE.

- При успешном выполнении запроса сервер должен вернуть статус 200 OK.

- В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.

Во всех обработчиках тип контента Content-Type — application/json.
*/
func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	// Значение параметра из URL
	id := chi.URLParam(r, "id")

	// Поиск задачи по значению параметра
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// Удаление задачи
	delete(tasks, id)

	// Тип контента и статус
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// Регистрация обработчиков
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTaskById)
	r.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
