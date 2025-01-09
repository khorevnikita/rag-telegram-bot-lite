package services

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"time"
)

func DrawLine(startTime time.Time, data []DateValues, title string, xLabel string, yLabel string) string {
	// Создаем карту для быстрого доступа к значениям по датам
	dataMap := make(map[string]int)
	for _, d := range data {
		dataMap[d.Date.Format("2006-01-02")] = d.Count
	}

	// Генерируем данные для каждой даты от startTime до сегодняшнего дня
	endTime := time.Now()
	var points plotter.XYs
	for t := startTime; !t.After(endTime); t = t.AddDate(0, 0, 1) {
		formattedDate := t.Format("2006-01-02")
		count := 0
		if val, exists := dataMap[formattedDate]; exists {
			count = val
		}

		points = append(points, plotter.XY{
			X: float64(t.Unix()),
			Y: float64(count),
		})
	}

	// Создание графика
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	// Настройка оси Y
	p.Y.Min = 0

	// Добавление линии
	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	line.Color = color.RGBA{R: 255, A: 0, B: 0}
	p.Add(line)

	// Форматирование оси X (даты)
	p.X.Tick.Marker = plot.TimeTicks{Format: "02.01.2006"}

	// Сохранение в PNG
	filePath := fmt.Sprintf("storage/%s.png", title)
	err = p.Save(8*vg.Inch, 4*vg.Inch, filePath)
	if err != nil {
		panic(err)
	}

	return filePath
}
