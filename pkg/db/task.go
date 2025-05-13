package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}
func Tasks(limit int) ([]*Task, error) {

	query := `SELECT id, date, title, comment, repeat FROM scheduler`

	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса SELECT: %w", err)
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		var id int64

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки результата: %w", err)
		}

		task.ID = fmt.Sprintf("%d", id)

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после итерации по строкам результата: %w", err)
	}

	return tasks, nil
}
func GetTask(id string) (*Task, error) {

	var taskID int
	_, err := fmt.Sscan(id, &taskID)
	if err != nil {
		return nil, fmt.Errorf("неверный формат ID: %w", err)
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	task := &Task{}

	err = db.QueryRow(query, taskID).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("ошибка базы данных при получении задачи: %w", err)
	}

	return task, nil
}

func UpdateTask(task *Task) error {

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ?  WHERE id = ? `

	res, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)

	if err != nil {

		return err
	}

	count, err := res.RowsAffected()
	if err != nil {

		return err
	}

	if count == 0 {

		return fmt.Errorf(`incorrect id for updating task`)

	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
