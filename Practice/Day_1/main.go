/*
	Используя заготовку из Lectio3_BooksAPI, добавить Rest API:

	PUT /books/{id}/  // Изменение существующей книги по id (id при этом не меняем!)
	DELETE /books/{id}// Удаление книги по id
	DELETE /books/    // Удаление всех книг

*/

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

func deleteAllBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	w.Header().Set("Content-Type", "application/json")
	books = nil
	json.NewEncoder(w).Encode(books)
	fmt.Fprintf(w, "<h1>All books has been deleted</h1>") //надо прописать в json формате, а не в html
}

func deleteBookByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	for _, book := range books {
		if book.ID == id {
			//delete(book, "{ID}")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
	json.NewEncoder(w).Encode(books)
}

/*
func getBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("now handling URL Path:", r.URL.Path)
	id := r.PathValue("id")
	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}
*/

/*
func createBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("now handling URL Path:", r.URL.Path)
	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	books = append(books, newBook)
	json.NewEncoder(w).Encode(newBook)
}
*/

func main() {

	//явно создаём router
	mux := http.NewServeMux() //Multiplexer

	// Добавляем маршруты к роутеру
	mux.HandleFunc("DELETE /books/", deleteAllBooks)
	mux.HandleFunc("DELETE /books/{id}", deleteBookByID)
	mux.HandleFunc("GET /books/", getBooks)
	//mux.HandleFunc("POST /books/", createBooks)

	fmt.Println("Start local web server...")

	// Передаём в ListenAndServe наш роутер (mux)
	log.Fatal(http.ListenAndServe("localhost"+port, mux)) // Запускаем веб сервер для "слушанья" (Завёрнутый в отслеживание ошибки)
}
