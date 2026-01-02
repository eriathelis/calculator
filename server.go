package main

import (
	"context"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"
)

// Шаблон загружаем один раз при старте
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// Структура для передачи данных в шаблон
type PageData struct {
	Result string
	Error  string
}

// Глобальная переменная для сервера
var srv *http.Server

// Основной обработчик формы калькулятора
func handler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")
	op := r.URL.Query().Get("op")

	data := PageData{}

	if aStr != "" && bStr != "" && op != "" {
		a, err1 := strconv.ParseFloat(aStr, 64)
		b, err2 := strconv.ParseFloat(bStr, 64)

		if err1 != nil || err2 != nil {
			data.Error = "Некорректные числа"
		} else {
			result, err := calculate(a, b, op)
			if err != nil {
				data.Error = err.Error()
			} else {
				if result == math.Trunc(result) {
					data.Result = fmt.Sprintf("%.0f", result)
				} else {
					data.Result = fmt.Sprintf("%.2f", result)
				}
			}
		}
	}

	tmpl.Execute(w, data)
}

// Обработчик кнопки "Выход"
func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, `
  <!DOCTYPE html>
  <html lang="ru">
  <head><meta charset="UTF-8"><title>Сервер остановлен</title></head>
  <body>
   <h2>Сервер остановлен/.</h2>
   <p>Можно закрыть вкладку браузера.</p>
  </body>
  </html>
 `)

	go func() {
		time.Sleep(500 * time.Millisecond)
		if err := srv.Shutdown(context.Background()); err != nil {
			fmt.Println("Ошибка при остановке сервера:", err)
		}
	}()
}

// Запуск сервера
func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/shutdown", shutdownHandler)

	srv = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Ошибка сервера:", err)
	}
}
