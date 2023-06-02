package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/abhishekmaurya0/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntries(t *testing.T) Entries {
	args := CreateEntriesParams{
		AccountID: CreateRandomAccount(t).ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := TestQuaries.CreateEntries(context.TODO(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.AccountID, entry.AccountID)
	return entry
}
func TestCreateEntries(t *testing.T) {
	CreateRandomEntries(t)
}
func TestGetEntries(t *testing.T) {
	entry1 := CreateRandomEntries(t)
	entry2, err := TestQuaries.GetEntries(context.TODO(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}
func TestDeleteEntries(t *testing.T) {
	entry1 := CreateRandomEntries(t)
	err := TestQuaries.DeleteEntries(context.TODO(), entry1.ID)
	require.NoError(t, err)
	entry2, err := TestQuaries.GetEntries(context.TODO(), entry1.ID)
	require.Error(t, err)
	require.Empty(t, entry2)
	require.Error(t, err, sql.ErrNoRows)
}
func TestListEntries(t *testing.T) {
	n := util.RandomInt(0, 101)
	for i := 0; i < int(n); i++ {
		_ = CreateRandomEntries(t)
	}
	lim := util.RandomInt(1, n)
	args := ListEntriesParams{
		Offset: int32(util.RandomInt(1, n)),
		Limit:  int32(lim),
	}
	entries, err := TestQuaries.ListEntries(context.TODO(), args)
	require.NoError(t, err)
	require.Len(t, entries, int(lim))
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
