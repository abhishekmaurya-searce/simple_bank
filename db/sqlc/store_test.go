package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store := NewStore(TestDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println("acc1 bal: ", account1.Balance, "acc2 bal: ", account2.Balance)
	n := 10
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.TODO(), TransferTxParam{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	//check for error and result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		// checking transefer
		transfer := result.Tranfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccount, account1.ID)
		require.Equal(t, transfer.ToAccount, account2.ID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransactions(context.TODO(), transfer.ID)
		require.NoError(t, err)

		//Check Entries
		//To Entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntries(context.TODO(), toEntry.ID)
		require.NoError(t, err)
		//From Entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntries(context.TODO(), fromEntry.ID)
		require.NoError(t, err)
		//check from account balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		//check to account balance
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// Check Account Balance (To Be done)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
		fmt.Println("acc1 bal: ", fromAccount.Balance, "acc2 bal: ", toAccount.Balance)
	}
	//check the updated accounts
	updatedAccount1, err := TestQuaries.GetAccounts(context.TODO(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := TestQuaries.GetAccounts(context.TODO(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, updatedAccount1.Balance, account1.Balance-int64(n)*amount)
	require.Equal(t, updatedAccount2.Balance, account2.Balance+int64(n)*amount)
	fmt.Println("acc1 bal: ", account1.Balance, "acc2 bal: ", account2.Balance)
}
func TestStoreDeadLock(t *testing.T) {
	store := NewStore(TestDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println("acc1 bal: ", account1.Balance, "acc2 bal: ", account2.Balance)
	n := 10
	amount := int64(10)
	errs := make(chan error)
	//results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		FromAccountId := account1.ID
		ToAccountId := account2.ID
		if i%2 == 1 {
			FromAccountId = account2.ID
			ToAccountId = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.TODO(), TransferTxParam{
				FromAccountId: FromAccountId,
				ToAccountId:   ToAccountId,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	//check for error and result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	//check the updated accounts
	updatedAccount1, err := TestQuaries.GetAccounts(context.TODO(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := TestQuaries.GetAccounts(context.TODO(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, updatedAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)
	fmt.Println("acc1 bal: ", account1.Balance, "acc2 bal: ", account2.Balance)
}
