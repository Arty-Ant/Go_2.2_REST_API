package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port = ":3000"

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func init() {
	books = []Book{
		{ID: "1", Title: "Приключения Тома Соейра", Author: "Марк Твен"},
		{ID: "2", Title: "Война и мир", Author: "Лев Толстой"},
		{ID: "3", Title: "Programming C", Author: "Брайн Керниганб Денис Ритчи"},
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func main() {

	//явно создаём router
	mux := http.NewServeMux()
	mux.HandleFunc("GET /books/", getBooks)

	fmt.Println("Start local web server...")
	log.Fatal(http.ListenAndServe("localhost"+port, mux)) // Запускаем веб сервер для "слушанья" (Завёрнутый в отслеживание ошибки)
}
