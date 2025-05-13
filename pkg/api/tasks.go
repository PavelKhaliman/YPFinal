package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go1f/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

		writeJSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("метод %s не разрешен", r.Method))
		return
	}

	//подсмотрел как сделать сортировку надеюсь таким методом это реализовывать правильно

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "asc"
	}
	tasks, err := db.Tasks(50, sort) // в параметре максимальное количество записей
	if err != nil {
		// здесь вызываете функцию, которая возвращает ошибку в JSON
		// её желательно было реализовать на предыдущем шаге
		// ...
		writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("ошибка при получении задач из базы данных: %w", err))
		return

	}
	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJSONResponse(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}

// GET
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	//ID из запроса
	id := r.URL.Query().Get("id")
	if id == "" {

		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("нет ID задачи"))
		return
	}

	//GET из DB
	task, err := db.GetTask(id)
	if err != nil {

		if err == sql.ErrNoRows {

			writeJSONError(w, http.StatusNotFound, fmt.Errorf("задача не найдена"))
		} else if err.Error() == fmt.Sprintf("неверный формат ID: %s", id) {
			writeJSONError(w, http.StatusBadRequest, fmt.Errorf("неверный формат идентификатора задачи"))
		} else {

			writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("ошибка при получении задачи: %w", err))
		}
		return
	}

	writeJSONResponse(w, http.StatusOK, task)
}

// PUT такойже как и POST поменять ADD на Update
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("ошибка десериализации JSON: %w", err))
		return
	}
	defer r.Body.Close()

	if task.Title == "" {
		writeJSONError(w, http.StatusBadRequest, errors.New("пустое поле"))
		return
	}

	err = checkDate(&task)
	if err != nil {

		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("ошибка добавления задачи в базу данных: %w", err))
		return
	}

	response := map[string]string{"id": task.ID}
	writeJSONResponse(w, http.StatusOK, response)
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("метод %s не разрешен", r.Method))
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("нет ID задачи"))
		return
	}

	// Получить задачу по id
	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, fmt.Errorf("задача не найдена"))
		} else {
			writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("ошибка при получении задачи: %w", err))
		}
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("ошибка при удалении задачи: %w", err))
			return
		}

		writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
		return
	}

	nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("ошибка вычисления следующей даты: %w", err))
		return
	}

	err = db.UpdateDate(nextDate, id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("ошибка при обновлении даты: %w", err))
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("метод %s не разрешен", r.Method))
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("нет ID задачи"))
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if err.Error() == "задача не найдена" {
			writeJSONError(w, http.StatusNotFound, err)
		} else {
			writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("ошибка при удалении задачи: %w", err))
		}
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
}
