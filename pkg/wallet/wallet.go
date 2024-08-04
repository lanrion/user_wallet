package wallet

import (
	"errors"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"sync"
)

type Operation int

// nolint
func (op Operation) String() string {
	switch op {
	case Deposit:
		return "Deposit"
	case Withdraw:
		return "Withdraw"
	case Transfer:
		return "Transfer"
	}
	return "Unknown"
}

const (
	Deposit  Operation = 0
	Withdraw Operation = 5
	Transfer Operation = 10
)

var balanceIncrease = math.Pow10(5)

var BucketManager *Bucket

type Bucket struct {
	lock    sync.RWMutex
	wallets map[uint32]*Wallet
}

func NewBucket(data map[uint32]*Wallet) *Bucket {
	if len(data) == 0 {
		data = make(map[uint32]*Wallet, 10000)
	}
	return &Bucket{wallets: data}
}

func (bucket *Bucket) Fetch(uid uint32) *Wallet {
	bucket.lock.RLock()
	if wallet, ok := bucket.wallets[uid]; ok {
		bucket.lock.RUnlock()
		return wallet
	}
	bucket.lock.RUnlock()

	bucket.lock.Lock()
	wallet := NewWallet(uid, 0)
	bucket.wallets[uid] = wallet
	bucket.lock.Unlock()
	return wallet
}

type Wallet struct {
	uid     uint32
	balance int64 // todo atomic
}

func NewWallet(uid uint32, balance float64) *Wallet {
	return &Wallet{uid: uid, balance: multiple8(balance)}
}

// Deposit 加钱
func (w *Wallet) Deposit(amount float64) error {
	log.Printf("[Wallet]Deposit uid: %d, amount: %f", w.uid, amount)
	w.balance = w.balance + multiple8(amount)
	return nil
}

// Withdraw 提现
func (w *Wallet) Withdraw(amount float64) error {
	log.Printf("[Wallet]Withdraw uid: %d, amount: %v, oldBalance: %v", w.uid, amount, w.balance)
	deductTotal := multiple8(amount)
	if w.balance < deductTotal {
		return errors.New("not enough balance")
	}
	w.balance = w.balance - deductTotal
	return nil
}

// Transfer TODO rollback...try
func (w *Wallet) Transfer(toUID uint32, amount float64) error {
	log.Printf("[Wallet]Transfer uid: %d, amount: %f, oldBalance: %v", toUID, amount, w.balance)
	err := w.Withdraw(amount)
	if err != nil {
		return err
	}
	toWallet := BucketManager.Fetch(toUID)
	if err = toWallet.Deposit(amount); err != nil {
		if err := w.Deposit(amount); err != nil {
			panic(err)
		}
	}
	return nil
}

func multiple8(amount float64) int64 {
	return int64(amount * balanceIncrease)
}

// Balance 返回余额
func (w *Wallet) Balance() float64 {
	if w.balance == 0 {
		return 0
	}
	d1 := decimal.NewFromInt(w.balance)
	d2 := decimal.NewFromFloat(balanceIncrease)
	return d1.Div(d2).InexactFloat64()
}
