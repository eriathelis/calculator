package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

// Шаблон загружаем один раз при старте (если файл отсутствует — программа упадёт и мы увидим ошибку сразу)
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// Структура данных, которую мы передаём в шаблон
type PageData struct {
	Result string
	Error  string
}

// handler обрабатывает запросы к корню сайта "/"
func handler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL: ?a=...&op=...&b=...
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")
	op := r.URL.Query().Get("op")

	// Подготовим структуру для шаблона
	data := PageData{}

	// Если какие-то параметры не заданы — просто отрисуем страницу без результата
	if aStr == "" || bStr == "" || op == "" {
		tmpl.Execute(w, data)
		return
	}

	// Пытаемся превратить строки в числа
	a, err1 := strconv.ParseFloat(aStr, 64)
	b, err2 := strconv.ParseFloat(bStr, 64)
	if err1 != nil || err2 != nil {
		data.Error = "Некорректные числа"
		tmpl.Execute(w, data)
		return
	}

	// Вызываем логику калькулятора (в calc.go)
	result, err := calculate(a, b, op)
	if err != nil {
		data.Error = err.Error()
	} else {
		// Красивый вывод: если число целое — без дробной части, иначе — с двумя знаками
		if result == math.Trunc(result) {
			data.Result = fmt.Sprintf("%.0f", result)
		} else {
			data.Result = fmt.Sprintf("%.2f", result)
		}
	}

	// Отдаём шаблон с данными
	tmpl.Execute(w, data)
}

// Функция для запуска сервера (вызывается из main.go)
func StartServer() {
	http.HandleFunc("/", handler)
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
