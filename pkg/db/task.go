package db

import (
	"fmt"
	"time"
)

const (
	storageDateFormat = "20060102"
	inputDateFormat   = "02.01.2006"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	insertQuery := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`
	res, err := db.Exec(insertQuery, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func Tasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	var getQuery string
	var args []any

	// Выбираем запрос по содержимому search: без поиска возвращаем ближайшие задачи,
	// дату ищем по полю date, остальные значения — по title и comment
	if search == "" {
		args = append(args, limit)
		getQuery = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?
	`
	} else if date, err := time.Parse(inputDateFormat, search); err != nil {
		pattern := "%" + search + "%"
		args = append(args, pattern, pattern, limit)
		getQuery = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ?
		ORDER BY date
		LIMIT ?
	`
	} else {
		args = append(args, date.Format(storageDateFormat), limit)
		getQuery = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE date = ?
		LIMIT ?
	`
	}

	rows, err := db.Query(getQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}

		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	task := &Task{}

	getQuery := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`
	if err := db.QueryRow(getQuery, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return nil, err
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	updateQuery := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`
	res, err := db.Exec(updateQuery, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func DeleteTask(id string) error {
	deleteQuery := `
		DELETE FROM scheduler
		WHERE id = ?
	`
	res, err := db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for deleting task")
	}
	return nil
}

func UpdateDate(next, id string) error {
	updateDateQuery := `
		UPDATE scheduler
		SET date = ?
		WHERE id = ?
	`
	res, err := db.Exec(updateDateQuery, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating date")
	}
	return nil
}
