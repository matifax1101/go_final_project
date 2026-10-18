package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrTaskNotFound = errors.New("задача не найдена")

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	res, err := database.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := database.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC, id ASC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*Task, 0)
	for rows.Next() {
		t := &Task{}
		var id int64
		if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = fmt.Sprint(id)
		result = append(result, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func GetTask(id string) (*Task, error) {
	var t Task
	var numericID int64
	err := database.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&numericID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	t.ID = fmt.Sprint(numericID)
	return &t, nil
}

func UpdateTask(task *Task) error {
	res, err := database.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func UpdateDate(next, id string) error {
	res, err := database.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := database.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}
