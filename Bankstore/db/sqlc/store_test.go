package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB) // testDB == nil -> store == nil
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("Before:", account1.Balance, account2.Balance)
	n := 8
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	// run n concurrent transfer transactions
	for i := range n {
		txName := fmt.Sprintf("tx %d", i+1) // "tx 1", 'tx 2", "tx 3"
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	// Check results
	existed := make(map[int]bool)
	for range n {
		// check errors
		err := <-errs
		require.NoError(t, err)
		// check transfers
		result := <-results
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		// Additional check
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		// check account's balance
		fmt.Println(">> Tx", account1.Balance, account2.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	// После отработки n горутин,проверяем результаты работы(в первую очередь, балансы)
	// final check
	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println("After:", updateAccount1.Balance, updateAccount2.Balance)
	// Добавляем все списания
	require.Equal(t, account1.Balance, updateAccount1.Balance+amount*int64(n))
	// Вычитаем все начисления
	require.Equal(t, account2.Balance, updateAccount2.Balance-amount*int64(n))
}
