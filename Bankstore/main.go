package main

import (
	"Bankstore/api"
	db "Bankstore/db/sqlc"
	"Bankstore/utils"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Загружаем настройки из файла app.env
	config, err := utils.LoadConfig(".") // "." - current directory
	if err != nil {
		log.Fatal("can not read config file", err)
	}
	// Соединение с БД
	pool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}
	// Закрываем соединение
	defer pool.Close()

	// Создаём необходимые экземпляры для работы
	store := db.NewStore(pool)     // хранилище
	server := api.NewServer(store) // роутинг и прочее

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can not stert server", err)
	}
}
