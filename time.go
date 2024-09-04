package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TimePicker struct {
	widget.BaseWidget
	hour      int
	minute    int
	second    int
	onChanged func(hour, minute, second int)
}

func NewTimePicker(initialTime time.Time, onChanged func(hour, minute, second int)) *TimePicker {
	tp := &TimePicker{
		hour:      initialTime.Hour(),
		minute:    initialTime.Minute(),
		second:    initialTime.Second(),
		onChanged: onChanged,
	}
	tp.ExtendBaseWidget(tp)
	return tp
}

func (tp *TimePicker) CreateRenderer() fyne.WidgetRenderer {
	hourSelect := widget.NewSelect(generateOptions(0, 23), func(s string) {
		tp.hour, _ = strconv.Atoi(s)
		tp.onChanged(tp.hour, tp.minute, tp.second)
	})
	hourSelect.SetSelected(fmt.Sprintf("%02d", tp.hour))

	minuteSelect := widget.NewSelect(generateOptions(0, 59), func(s string) {
		tp.minute, _ = strconv.Atoi(s)
		tp.onChanged(tp.hour, tp.minute, tp.second)
	})
	minuteSelect.SetSelected(fmt.Sprintf("%02d", tp.minute))

	secondSelect := widget.NewSelect(generateOptions(0, 59), func(s string) {
		tp.second, _ = strconv.Atoi(s)
		tp.onChanged(tp.hour, tp.minute, tp.second)
	})
	secondSelect.SetSelected(fmt.Sprintf("%02d", tp.second))

	content := container.NewHBox(
		hourSelect,
		widget.NewLabel(":"),
		minuteSelect,
		widget.NewLabel(":"),
		secondSelect,
	)

	return widget.NewSimpleRenderer(content)
}

func generateOptions(start, end int) []string {
	options := make([]string, end-start+1)
	for i := start; i <= end; i++ {
		options[i-start] = fmt.Sprintf("%02d", i)
	}
	return options
}
