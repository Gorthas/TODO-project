package db

import (
	"time"
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

	if search == "" {
		args = append(args, limit)
		getQuery = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?
	`
	} else if date, err := time.Parse("02.01.2006", search); err != nil {
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
		args = append(args, date.Format("20060102"), limit)
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
