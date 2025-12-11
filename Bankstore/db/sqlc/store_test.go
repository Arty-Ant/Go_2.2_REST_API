package db

import (
	"context"
	"fmt"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("Before:", account1.Balance, account2.Balance)

	amount := int64(10)

	result, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result: %v\n", result)
}
