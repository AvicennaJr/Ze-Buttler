package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"fyne.io/fyne/v2/widget"
)

func imageToBytes(imagePath string) ([]byte, error) {

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error opening image file:", err)
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil

}

func refreshTodoList(db *sql.DB, list *widget.Table) error {
	err := deletePastDueTodos(db)
	if err != nil {
		return err
	}

	todos, err := listTodos(db)
	if err != nil {
		return err
	}

	listTodoView = [][]string{{"ID", "Title", "Deadline"}}
	for _, todo := range todos {
		listTodoView = append(listTodoView, []string{
			fmt.Sprintf("%d", todo.ID),
			todo.Title,
			todo.Deadline,
		})
	}

	list.Refresh()
	return nil
}
