/*
	Используя заготовку из Lection3_BooksAPI, добавить API к ресурсу books:

	PUT /books/{id}/   // Изменение существующей книги по id(id при этом не меняем!)
	DELETE /books/{id} // Удаление книги по id
	DELETE /books/     // Удаление всех книг
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices" // там почти всё есть: Insert, Contains, Delete, Replace, Sort, Reverse
)

const port = ":3000"

type Book struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
}

var books []Book


type MessageResponse struct {
	Message string `json:"message"`
}


func init() {
	books = []Book{
		{ID: "1", Title: "Приключения Тома Соейра", Author: "Марк Твен"},
		{ID: "2", Title: "Война и мир", Author: "Лев Толстой"},
		{ID: "3", Title: "Programming C", Author: "Брайн Керниган, Денис Ритчи"},

	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("Method: %s, handling URL Path: %s\n", r.Method, r.URL.Path)
	json.NewEncoder(w).Encode(books)
}

func getBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("Method: %s, handling URL Path: %s\n", r.Method, r.URL.Path)

	id := r.PathValue("id")
	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, fmt.Sprintf("Book with id=%s not found", id), http.StatusNotFound)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("Method: %s, handling URL Path: %s\n", r.Method, r.URL.Path)

	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	books = append(books, newBook)
	json.NewEncoder(w).Encode(newBook)
}

func deleteBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("Method: %s, handling URL Path: %s\n", r.Method, r.URL.Path)

	// Variant 1
	// clear(books) // очистка данных, но не удаление самых экземпляров
	// w.WriteHeader(http.StatusNoContent) // Правильный результат работы, 204 status code(no content)

	// Variant 2
	books = slices.Delete(books, 0, len(books)) 
	
	message := MessageResponse{
		Message: "All books have deleted",
	}
	json.NewEncoder(w).Encode(message)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("Method: %s, handling URL Path: %s\n", r.Method, r.URL.Path)


	var modifiedBook Book
	err := json.NewDecoder(r.Body).Decode(&modifiedBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var bookIndex int = -1
	id := r.PathValue("id")
	for index, book := range books {
		if book.ID == id {
			bookIndex = index
			break
		}
	}

	if bookIndex == -1 {
		http.Error(w, fmt.Sprintf("Book with id=%s not found", id), http.StatusNotFound)
		return 
	}

	modifiedBook.ID = id
	books[bookIndex] = modifiedBook

	json.NewEncoder(w).Encode(modifiedBook)
}

func deleteBookById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("Method: %s, handling URL Path: %s\n", r.Method, r.URL.Path)

	var indexBook int = -1
	id := r.PathValue("id")

	for index, book := range books {
		if book.ID == id {
			indexBook = index
			break
		}
	}

	if indexBook == -1 {
		http.Error(w, fmt.Sprintf("Book with id=%s not found", id), http.StatusNotFound)
		return
	}
	// Variant 1
	// books = append(books[:indexBook], books[indexBook+1:]...)

	// Variant 2
	books = slices.Delete(books, indexBook, indexBook+1)

	message := MessageResponse{
		Message: fmt.Sprintf("Book with id=%s has deleted", id),
	}
	json.NewEncoder(w).Encode(message)
}


func main() {
	// явно создаём router
	mux := http.NewServeMux() // Multiplexer

	// Добавляем маршруты к роутеру
	mux.HandleFunc("GET /books/", getBooks) 
	mux.HandleFunc("GET /books/{id}/", getBookById) 
	mux.HandleFunc("POST /books/", createBook) 
	mux.HandleFunc("PUT /books/{id}/", updateBook)
	mux.HandleFunc("DELETE /books/", deleteBooks)
	mux.HandleFunc("DELETE /books/{id}/", deleteBookById)
	fmt.Println("Start local web server...")

	// Передаем в ListenAndServer наш роутер(mux)
	log.Fatal(http.ListenAndServe("localhost"+port, mux)) 
}