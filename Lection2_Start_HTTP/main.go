package main

import (
	"fmt"
	"log"
	"net/http"
)

// w- responceWrite (куда записывать ответ)
// r - request (откуда брать запрос)

// Обработчик
func GetGreet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello, I'm new web srver! </h1>")
	fmt.Println("Method:", r.Method)
	fmt.Println("URL", r.URL)
}

func main() {
	http.HandleFunc("GET /", GetGreet)           // Если придёт запрос по адресу (ресурсу) "/"(root), то вызывается GetGreet
	log.Fatal(http.ListenAndServe(":8080", nil)) // Запускаем веб сервер для "слушанья" (Завёрнутый в отслеживание ошибки)
}
