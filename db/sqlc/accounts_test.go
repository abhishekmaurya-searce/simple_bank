package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/abhishekmaurya0/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Accounts {
	args := CreateAccountsParams{
		Owner:    util.RandomOwner(), //Randomly genarated
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := TestQuaries.CreateAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccounts(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := TestQuaries.GetAccounts(context.TODO(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccounnts(t *testing.T) {
	account1 := CreateRandomAccount(t)
	args := UpdateAccountsParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	account2, err := TestQuaries.UpdateAccounts(context.TODO(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, args.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccounts(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := TestQuaries.DeleteAccounts(context.TODO(), account1.ID)
	require.NoError(t, err)
	account2, err := TestQuaries.GetAccounts(context.TODO(), account1.ID)
	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows)
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	n := int64(10)
	for i := 0; i < int(n); i++ {
		_ = CreateRandomAccount(t)
	}
	lim := int32(util.RandomInt(1, n))
	args := ListAccountsParams{
		Offset: int32(util.RandomInt(1, n)),
		Limit:  lim,
	}
	accounts, err := TestQuaries.ListAccounts(context.TODO(), args)
	require.NoError(t, err)
	require.Len(t, accounts, int(lim))
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
