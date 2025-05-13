package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"go1f/pkg/db"
	"net/http"
	"strconv"
	"time"
)

// п.3 из задания создать новую функцию checkDate
func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("неверная дата %s: %w", task.Date, err)
	}

	//валидность даты
	year, month, day := t.Date()
	reconstructed := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	if reconstructed.Year() != year || reconstructed.Month() != month || reconstructed.Day() != day {
		return fmt.Errorf("некорректная дата: %s", task.Date)
	}

	// сравниваем только по дню
	tDate := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if tDate.Before(nowDate) {

		if len(task.Repeat) == 0 {

			task.Date = nowDate.Format(DateFormat)
		} else {

			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка в правиле повторения: %w", err)
			}
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	//1.Нужно десериализовать полученный в запросе JSON
	var task db.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Errorf("ошибка десериализации JSON: %w", err))
		return
	}
	defer r.Body.Close()

	// 2. Проверить, что поле task.Title не пустое.
	if task.Title == "" {
		writeJSONError(w, http.StatusBadRequest, errors.New("пустое поле"))
		return
	}
	// 3. Проверить на корректность полученное значение task.Date. Это лучше сделать в отдельной функции checkDate(task *db.Task) error

	err = checkDate(&task)
	if err != nil {

		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	// 4. Пришла очередь вызвать функцию db.AddTask(task), чтобы добавить задачу в базу данных.
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("ошибка добавления задачи в базу данных: %w", err))
		return
	}

	// 5. Осталось возвратить идентификатор добавленной задачи в виде JSON.
	response := map[string]string{"id": strconv.FormatInt(id, 10)}
	writeJSONResponse(w, http.StatusOK, response)
}

// writeJSONResponse записывает JSON ответ с указанным статусом.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		fmt.Printf("Ошибка при записи JSON ответа: %v\n", err)
	}
}

// writeJSONError записывает JSON ответ с ошибкой и указанным статусом.
func writeJSONError(w http.ResponseWriter, statusCode int, err error) {
	writeJSONResponse(w, statusCode, map[string]string{"error": err.Error()})
}
