/*
	Pattern: <filename>_test.go
*/

package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"Bankstore/utils"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
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
	createRandomAccount(t)
}

func createRandomAccount(t *testing.T) Account {
	ra := utils.RandomAccount()
	arg := CreateAccountParams{
		Owner:   ra.Owner,
		Balance: utils.RandomInt(0, 1000),
		//Balance:  ra.Balance,
		Currency: Currency(ra.Currency),
	}
	account, err := testQueries.CreateAccount(ctx, arg)

	// Две проверки на результат работы CreateAccount
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestGetAccount(t *testing.T) {

	// Создание тестового аккаунта
	acc1 := createRandomAccount(t)

	// Вызов тестируемого метода
	acc2, err := testQueries.GetAccount(ctx, acc1.ID)

	// Проверки
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)

	require.WithinDuration(t, acc1.CreatedAt.Time, acc2.CreatedAt.Time, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	// Delete created account
	err := testQueries.DeleteAccount(ctx, account1.ID)
	require.NoError(t, err)
	account2, err := testQueries.GetAccount(ctx, account1.ID)
	require.Error(t, err)
	require.Empty(t, account2)
	// проверка на тип ошибки(нуль строк было обработано)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
}
func TestListAccounts(t *testing.T) {
	for range 10 {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  8,
		Offset: 15,
	}
	accounts, err := testQueries.ListAccounts(ctx, arg)
	require.NoError(t, err)
	require.Len(t, accounts, 8)
	for _, acc := range accounts {
		require.NotZero(t, acc)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomInt(0, 1000),
	}
	account2, err := testQueries.UpdateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	// test all fields
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}
