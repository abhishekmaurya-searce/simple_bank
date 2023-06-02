package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}

}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {

			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// Contains Input parameter for Money Transfer Transaction
type TransferTxParam struct {
	ToAccountId   int64 `json:"to_account_id"`
	FromAccountId int64 `json:"from_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Tranfer     Transactions `json:"transfer"`
	FromAccount Accounts     `json:"from_account"`
	ToAccount   Accounts     `json:"to_account"`
	FromEntry   Entries      `json:"from_entry"`
	ToEntry     Entries      `json:"to_entry"`
}

func (s *Store) TransferTx(ctx context.Context, args TransferTxParam) (TransferTxResult, error) {

	var result TransferTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Tranfer, err = q.CreateTransactions(ctx, CreateTransactionsParams{
			FromAccount: args.FromAccountId,
			ToAccount:   args.ToAccountId,
			Amount:      args.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: args.FromAccountId,
			Amount:    -1 * args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: args.ToAccountId,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}
		// Todo the balance update as it requires locking cause have chances of deadlock
		if args.FromAccountId < args.ToAccountId {

			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, args.FromAccountId, -args.Amount, args.ToAccountId, args.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, args.ToAccountId, args.Amount, args.FromAccountId, -args.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return result, err
}
func AddMoney(
	ctx context.Context,
	q *Queries,
	account1ID int64,
	amount1 int64,
	account2ID int64,
	amount2 int64,
) (Account1 Accounts, Account2 Accounts, err error) {
	Account1, err = q.AddAccountsBalance(ctx, AddAccountsBalanceParams{
		ID:     account1ID,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	Account2, err = q.AddAccountsBalance(ctx, AddAccountsBalanceParams{
		ID:     account2ID,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
