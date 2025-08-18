package db

import (
	"GO_TODO-list/pkg/constants"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT :limit`
	rows, err := DB.Query(query, sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

func SearchTasks(search string, limit int) ([]*Task, error) {
	const searchDateFormat = "02.01.2006"
	var (
		query     string
		queryArgs []interface{}
	)

	//формирование запроса по дате или строке в поиске
	date, err := time.Parse(searchDateFormat, search)
	if err == nil {
		query = `SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE date = :date 
		ORDER BY date ASC
		LIMIT :limit`
		queryArgs = append(queryArgs,
			sql.Named("date", date.Format(constants.DATE_FORMAT)),
			sql.Named("limit", limit))

	} else {
		query = `SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE title LIKE :search OR comment LIKE :search 
		ORDER BY date ASC
		LIMIT :limit`
		like := "%" + search + "%"
		queryArgs = append(queryArgs,
			sql.Named("search", like),
			sql.Named("limit", limit))
	}

	//выполнение запроса в БД
	rows, err := DB.Query(query, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	err := DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler 
			SET date = :date,
				title = :title,
				comment = :comment,
				repeat = :repeat 
				WHERE id = :id`
	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
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

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := DB.Exec(query,
		sql.Named("date", next),
		sql.Named("id", id),
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
	query := `DELETE FROM scheduler WHERE id = :id`

	res, err := DB.Exec(query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
