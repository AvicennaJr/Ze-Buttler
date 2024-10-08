package main

import (
	"database/sql"
)

func createTodo(db *sql.DB, todo Todo) (Todo, error) {
	tx, err := db.Begin()
	if err != nil {
		return Todo{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO todo(title, deadline) VALUES (?, ?) RETURNING id")
	if err != nil {
		return Todo{}, err
	}
	defer stmt.Close()

	var id int32
	err = stmt.QueryRow(todo.Title, todo.Deadline).Scan(&id)
	if err != nil {
		return Todo{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Todo{}, err
	}

	todo.ID = id
	return todo, nil
}

func readTodo(db *sql.DB, id int32) (Todo, error) {
	var todo Todo
	err := db.QueryRow("SELECT id, title, deadline FROM todo WHERE id = ?", id).Scan(&todo.ID, &todo.Title, &todo.Deadline)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func updateTodo(db *sql.DB, todo Todo) error {
	_, err := db.Exec("UPDATE todo SET title = ?, deadline = ? WHERE id = ?", todo.Title, todo.Deadline, todo.ID)
	return err
}

func deleteTodo(db *sql.DB, id int32) error {
	_, err := db.Exec("DELETE FROM todo WHERE id = ?", id)
	return err
}

func listTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query(`
        SELECT id, title, deadline 
        FROM todo 
        ORDER BY datetime(
            substr(deadline, 7, 4) || '-' || 
            substr(deadline, 4, 2) || '-' || 
            substr(deadline, 1, 2) || ' ' || 
            substr(deadline, 12)
        )
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Deadline); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func deletePastDueTodos(db *sql.DB) error {
	_, err := db.Exec(`
        DELETE FROM todo 
        WHERE datetime(
            substr(deadline, 7, 4) || '-' || 
            substr(deadline, 4, 2) || '-' || 
            substr(deadline, 1, 2) || ' ' || 
            substr(deadline, 12)
        ) < datetime('now', 'localtime')
    `)
	return err
}
