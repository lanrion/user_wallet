package wallet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	BucketManager = NewBucket(nil)
	m.Run()
}

func TestOperation_String(t *testing.T) {
	e := Deposit.String()
	assert.Equal(t, "Deposit", e)
	e = Operation(33333).String()
	assert.Equal(t, "Unknown", e)

	e = Transfer.String()
	assert.Equal(t, "Transfer", e)

	e = Withdraw.String()
	assert.Equal(t, "Withdraw", e)

}

func TestBucket_Fetch(t *testing.T) {
	wal := BucketManager.Fetch(12)
	assert.Equal(t, int(wal.balance), 0)
	wal.Deposit(10)
	wal2 := BucketManager.Fetch(12)
	assert.Equal(t, 10.0, wal2.Balance())

}

func TestWallet_Withdraw(t *testing.T) {
	t.Run("withdraw fail when balance is zero", func(t *testing.T) {
		wallet := NewWallet(1, 0)
		err := wallet.Withdraw(10)
		assert.NotNil(t, err)
		assert.Equal(t, 0, int(wallet.Balance()))
	})

	t.Run("withdraw all", func(t *testing.T) {
		wallet := NewWallet(1, 10)
		err := wallet.Withdraw(10)
		assert.Nil(t, err)
		assert.Equal(t, 0, int(wallet.Balance()))
	})

	t.Run("withdraw overflow", func(t *testing.T) {
		wallet := NewWallet(1, 10)

		err := wallet.Withdraw(100)
		assert.NotNil(t, err)
	})
}

func TestWallet_Balance(t *testing.T) {
	wallet := NewWallet(1, 1001221)
	assert.Equal(t, float64(1001221), wallet.Balance())
}

func TestWallet_Deposit(t *testing.T) {
	wallet := NewWallet(1, 0)
	wallet.Deposit(10)
	assert.Equal(t, 10.0, wallet.Balance())

	wallet.Deposit(20)
	assert.Equal(t, 30.0, wallet.Balance())

}
