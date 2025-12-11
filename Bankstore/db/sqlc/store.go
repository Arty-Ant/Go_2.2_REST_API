package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
	*Queries
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) exectTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx) //создаём экземпляр Queries на основе pgx.Tx
	err = fn(q)  // обработка созданного экземпляра

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

// Начниаем ревлизацию перевода через транзакцию

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amont"`
}

type TransferTxResult struct {
	Transfer      Transfer `json:"transfer"`
	FromAccountID int64    `json:"from_account_id"`
	ToAccountID   int64    `json:"to_account_id"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.exectTx(ctx, func(q *Queries) error {
		//Step 1. Create Transfer record
		var err error
		result.Transfer, err  = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		// Step 2. Create two Entry records
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//Step 3. Update account's balance
		//TODO
		return nil
	})
	return result, err
}
