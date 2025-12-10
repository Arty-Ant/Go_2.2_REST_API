/*
	Pattern: <filename>_test.go
*/

package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"Bankstore/utils"
)

const (
	dbSource = "postgresql://app_user:pswd@localhost:5432/bankdb?sslmode=disable"
)

// Глобальный контекст для работы с БД и тестами
var ctx = context.Background()

// Декларация переменной
var testQueries *Queries

func TestMain(m *testing.M) {
	// Соединение с БД
	conn, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}
	// Закрываем соединение
	defer conn.Close(ctx)

	// Вызываем конструктор для создания экземпляра типа данных Queries
	testQueries = New(conn)

	// Запускаем subtest(тесты) и итоговый код выполнения передаем в Exit()
	os.Exit(m.Run())
}

func TestCreateAccount(t *testing.T) {
	ra := utils.RandomAccount()
	arg := CreateAccountParams{
		Owner:    ra.Owner,
		//Balance: utils.RandomInt(0, 1000),
		Balance:  ra.Balance,
		Currency: Currency(ra.Currency),
	}
	account, err := testQueries.CreateAccount(ctx, arg)

	// Две проверки на результат работы CreateAccount
	require.NoError(t, err)
	require.NotEmpty(t, account)
}
