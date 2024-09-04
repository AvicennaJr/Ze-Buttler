package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID       int32
	Title    string
	Deadline string
}

func onSelected(t string) {}

func onTimeSelected(hour, minute, second int) {}

var listTodoView = [][]string{{"ID", "Title", "Deadline"}}

func main() {
	iconBytes, err := imageToBytes("./icons/icon.png")
	if err != nil {
		log.Fatal(err)
	}

	icon := fyne.NewStaticResource(
		"logo.png",
		iconBytes,
	)
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}

	aiApp := app.New()
	listWindow := aiApp.NewWindow("Your Todos")

	todos, err := listTodos(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, todo := range todos {
		listTodoView = append(listTodoView, []string{
			fmt.Sprintf("%d", todo.ID),
			todo.Title,
			todo.Deadline,
		})
	}

	list := widget.NewTable(
		func() (int, int) {
			return len(listTodoView), len(listTodoView[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(listTodoView[i.Row][i.Col])
		},
	)

	list.SetColumnWidth(0, 50)
	list.SetColumnWidth(1, 150)
	list.SetColumnWidth(3, 200)

	listWindow.SetContent(list)
	listWindow.SetMaster()
	listWindow.Resize(fyne.NewSize(600, 400))

	listWindow.SetCloseIntercept(func() {
		listWindow.Hide()
	})

	createWindow := aiApp.NewWindow("Create Todo")

	title := widget.NewEntry()
	startingDate := time.Now()
	deadlineDate := NewCalendar(startingDate.Format("02-01-2006"), onSelected)
	deadlineTime := NewTimePicker(startingDate, onTimeSelected)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: title},
			{Text: "Deadline", Widget: deadlineDate},
			{Text: "DeadlineTime", Widget: deadlineTime},
		},
		OnSubmit: func() {
			dl := fmt.Sprintf("%s %02d:%02d:%02d", deadlineDate.selectedDate, deadlineTime.hour, deadlineTime.minute, deadlineTime.second)

			todo := Todo{
				Title:    title.Text,
				Deadline: dl,
			}
			todo, err = createTodo(db, todo)

			if err != nil {
				fmt.Println(err)
			}
			listTodoView = append(listTodoView, []string{
				fmt.Sprintf("%d", todo.ID),
				todo.Title,
				todo.Deadline,
			})
			list.Refresh()
			title.SetText("")

			createWindow.Hide()
		},
	}

	createWindow.Resize(fyne.NewSize(500, 200))
	createWindow.SetContent(form)

	createWindow.SetCloseIntercept(func() {
		createWindow.Hide()
	})

	aiResponse, err := CallAI(db)
	if err != nil {
		aiResponse = err.Error()
	}

	err = Alert("Ze Buttler", aiResponse, "./icons/icon.png")

	if err != nil {
		log.Println(err)
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			aiResponse, err := CallAI(db)
			if err != nil {
				aiResponse = err.Error()
			}
			err = Alert("Ze Buttler", aiResponse, "./icons/icon.png")

			if err != nil {
				log.Println(err)
			}

		}
	}()

	if desk, ok := aiApp.(desktop.App); ok {
		c := fyne.NewMenuItem("Create Todo", func() {
			createWindow.Show()
		})
		l := fyne.NewMenuItem("List Todos", func() {
			listWindow.Show()
		})
		m := fyne.NewMenu("AI", c, l)
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(icon)
	}

	aiApp.Run()
}