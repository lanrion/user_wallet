package wallet

import "user_wallet/pkg/internal/kfkmodule"

type Job struct {
	doneChan      chan struct{}
	operation     Operation
	CurrentUserID uint32 // 当前任务归属的用户id
	ToUserID      uint32
	Amount        float64
	Err           error
}

func NewJob(userID uint32, amount float64, operation Operation, toUserID uint32) *Job {
	return &Job{
		CurrentUserID: userID,
		Amount:        amount,
		doneChan:      make(chan struct{}),
		operation:     operation,
		ToUserID:      toUserID,
	}
}

func (j *Job) Wait() <-chan struct{} {
	return j.doneChan
}

func (j *Job) Finish() {
	close(j.doneChan)
}

func (j *Job) Process(csm Consumer) {
	defer j.Finish()

	wallet := BucketManager.Fetch(j.CurrentUserID)
	switch j.operation {
	case Deposit:
		j.Err = wallet.Deposit(j.Amount)
	case Withdraw:
		j.Err = wallet.Withdraw(j.Amount)
	case Transfer:
		j.Err = wallet.Transfer(j.ToUserID, j.Amount)
		toWallet := BucketManager.Fetch(j.ToUserID)
		pd := kfkmodule.PushData{
			UID:     j.ToUserID,
			Balance: toWallet.balance,
		}
		csm.Produce(&pd)
	}
	if j.Err != nil {
		return
	}
	pd := kfkmodule.PushData{
		UID:     wallet.uid,
		Balance: wallet.balance,
	}
	csm.Produce(&pd)

}
