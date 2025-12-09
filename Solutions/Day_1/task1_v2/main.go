/*
Используя заготовку из Lection3_BooksAPI, добавить API к ресурсу books:
PUT /books/{id}/ // Изменение существующей книги по id(id при этом не меняем!)
DELETE /books/{id} // Удаление книги по id
DELETE /books/ // Удаление всех книг
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

)
const port = ":3001"
type Book struct {
	ID string `json:"id"`
    Title string `json:"title"`
    Author string `json:"author"`
}
var books []Book
func init() {
    books = []Book{
        {ID: "1", Title: "Приключения Тома Соейра", Author: "Марк Твен"},
        {ID: "2", Title: "Война и мир", Author: "Лев Толстой"},
        {ID: "3", Title: "Programming C", Author: "Брайн Керниган, Денис Ритчи"},
    }
}
func getBooks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(books)
}
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
    http.Error(w, fmt.Sprintf("Book with id=%s not found", id), http.StatusNotFound)
}
func createBook(w http.ResponseWriter, r *http.Request) {
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


func deleteBooks(w http.ResponseWriter, r *http.Request){
	if r.Method != "DELETE" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Println("now handling delete URL Path:", r.URL.Path)
	books = []Book{}

}

func deleteBook(w http.ResponseWriter, r *http.Request){
	if r.Method != "DELETE" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Println("now handling delete URL Path:", r.URL.Path)
	
    id := r.PathValue("id")
	if id != "" {
		idx, _ := strconv.Atoi(id)
		books = append(books[:idx-1], books[idx:]...)
	}
}

func updateBook(w http.ResponseWriter, r *http.Request){
	if r.Method != "PUT" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Println("now handling update URL Path:", r.URL.Path)
    var updateBook Book
    err := json.NewDecoder(r.Body).Decode(&updateBook)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    id := r.PathValue("id")
    idx, _ := strconv.Atoi(id)
    books[idx-1].Author = updateBook.Author
    books[idx-1].Title = updateBook.Title
}

func main() {
    // явно создаём router
    mux := http.NewServeMux() // Multiplexer
    // Добавляем маршруты к роутеру
    mux.HandleFunc("GET /books/", getBooks)
    mux.HandleFunc("GET /books/{id}/", getBookById)
    mux.HandleFunc("POST /books/", createBook)
	mux.HandleFunc("DELETE /books/", deleteBooks)
	mux.HandleFunc("DELETE /books/{id}/", deleteBook)
    mux.HandleFunc("PUT /books/{id}/", updateBook)

    
    fmt.Println("Start local web server...")
    // Передаем в ListenAndServer наш роутер(mux)
    log.Fatal(http.ListenAndServe("localhost"+port, mux))
}