/*
	Pattern: <filename>_test.go
*/

package db

import (
	"context"
	"testing"
	"time"

	"Bankstore/utils"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

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
	account, err := testQueries.CreateAccount(context.Background(), arg)

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
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

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
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
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
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 10)
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
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	// test all fields
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}
