package main

import (
	"Bankstore/api"
	db "Bankstore/db/sqlc"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource      = "postgresql://app_user:pswd@localhost:5432/bankdb?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// Соединение с БД
	pool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}
	// Закрываем соединение
	defer pool.Close()

	// Создаём необходимые экземпляры для работы
	store := db.NewStore(pool)     // хранилище
	server := api.NewServer(store) // роутинг и прочее

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("can not stert server", err)
	}
}
