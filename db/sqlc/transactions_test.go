package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/abhishekmaurya0/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransaction(t *testing.T) Transactions {
	accounts, _ := TestQuaries.ListAccounts(context.TODO(), ListAccountsParams{
		Offset: int32(util.RandomInt(0, 100)),
		Limit:  2,
	})
	args := CreateTransactionsParams{
		ToAccount:   accounts[0].ID,
		FromAccount: accounts[1].ID,
		Amount:      util.RandomMoney(),
	}
	transaction, err := TestQuaries.CreateTransactions(context.TODO(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)
	require.Equal(t, args.Amount, transaction.Amount)
	require.Equal(t, args.ToAccount, transaction.ToAccount)
	require.Equal(t, args.FromAccount, transaction.FromAccount)
	require.NotZero(t, transaction.ID)
	require.NotZero(t, transaction.CreatedAt)
	return transaction
}
func TestCreateTransactions(t *testing.T) {
	_ = CreateRandomTransaction(t)
}
func TestGetTransaction(t *testing.T) {

	transaction1 := CreateRandomTransaction(t)
	transaction2, err := TestQuaries.GetTransactions(context.TODO(), transaction1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transaction2)
	require.Equal(t, transaction1, transaction2)
}
func TestDeleteTransaction(t *testing.T) {
	transaction1 := CreateRandomTransaction(t)
	err := TestQuaries.DeleteTransactions(context.TODO(), transaction1.ID)
	require.NoError(t, err)
	transaction, err := TestQuaries.GetTransactions(context.TODO(), transaction1.ID)
	require.Error(t, err, sql.ErrNoRows)
	require.Empty(t, transaction)
}

func TestListTransaction(t *testing.T) {
	n := util.RandomInt(1, 101)
	for i := 0; i < int(n); i++ {
		_ = CreateRandomTransaction(t)
	}
	lim := util.RandomInt(0, n)
	args := ListTransactionsParams{
		Offset: int32(util.RandomInt(0, n)),
		Limit:  int32(lim),
	}
	transactions, err := TestQuaries.ListTransactions(context.TODO(), args)
	require.NoError(t, err)
	require.Len(t, transactions, int(lim))
	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)

	}
}
