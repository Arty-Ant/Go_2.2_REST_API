package db

import (
	"context"
	"log"
	"os"
	"testing"

	"Bankstore/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// const (
// 	dbSource = "postgresql://app_user:pswd@localhost:5432/bankdb?sslmode=disable"
// )

// Глобальный контекст для работы с БД и тестами
//var ctx = context.Background()

// Декларация переменной
var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	// Загружаем настройки из файла app.env
	config, err := utils.LoadConfig("../..") // "." - current directory
	if err != nil {
		log.Fatal("can not read config file", err)
	}
	// Соединение с БД
	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}
	// Закрываем соединение
	defer testDB.Close()

	// Вызываем конструктор для создания экземпляра типа данных Queries
	testQueries = New(testDB)

	// Запускаем subtest(тесты) и итоговый код выполнения передаем в Exit()
	os.Exit(m.Run())
}
