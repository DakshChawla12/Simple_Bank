package db

import (
	"context"
	"testing"
	"time"

	"github.com/DakshChawla/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	testCases := []struct {
		name      string
		arg       CreateAccountParams
		expectErr bool
	}{
		{
			name: "valid USD account",
			arg: CreateAccountParams{
				Owner:    "test1",
				Balance:  1000,
				Currency: "USD",
			},
			expectErr: false,
		},
		{
			name: "valid INR account",
			arg: CreateAccountParams{
				Owner:    "test2",
				Balance:  1000,
				Currency: "USD",
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			account, err := testQueries.CreateAccount(context.Background(), tc.arg)

			if tc.expectErr {
				require.Error(t, err)
				require.Empty(t, account)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, account)

			require.Equal(t, tc.arg.Owner, account.Owner)
			require.Equal(t, tc.arg.Balance, account.Balance)
			require.Equal(t, tc.arg.Currency, account.Currency)

			require.NotZero(t, account.ID)
			require.NotZero(t, account.CreatedAt)
		})
	}
}

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: "USD",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	return account
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestGetAccount_NotFound(t *testing.T) {
	account, err := testQueries.GetAccount(context.Background(), 999999)

	require.Error(t, err)
	require.Empty(t, account)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	newBalance := util.RandomMoney()

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: newBalance,
	}

	updated, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updated)

	require.Equal(t, account.ID, updated.ID)
	require.Equal(t, account.Owner, updated.Owner)
	require.Equal(t, newBalance, updated.Balance)
	require.Equal(t, account.Currency, updated.Currency)
}

func TestUpdateAccount_NotFound(t *testing.T) {
	arg := UpdateAccountParams{
		ID:      999999,
		Balance: 100,
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)

	require.Error(t, err)
	require.Empty(t, account)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	// Verify deletion
	deleted, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, deleted)
}

func TestDeleteAccount_NotFound(t *testing.T) {
	err := testQueries.DeleteAccount(context.Background(), 999999)
	require.NoError(t, err)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
