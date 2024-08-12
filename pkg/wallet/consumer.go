package wallet

import (
	"context"
	"errors"
	"log"
	"user_wallet/pkg/internal/kfkmodule"
)

const userIdIdx = 10000

type Consumer interface {
	Start() error
	PushBack(ctx context.Context, job *Job) error
	Produce(*kfkmodule.PushData)
}

type consumer struct {
	queue  map[uint32]chan *Job // TODO 可以再拆小一点
	pusher Pusher
}

func NewConsumer(ps Pusher) Consumer {
	queue := make(map[uint32]chan *Job)

	return &consumer{queue: queue, pusher: ps}
}

func (c *consumer) Start() error {
	for i := 0; i < userIdIdx; i++ {
		idxCh := make(chan *Job, 10000)
		c.queue[uint32(i)] = idxCh
		go func(ch chan *Job) {
			for {
				select {
				case jobData := <-ch:
					jobData.Process(c)
				}
			}
		}(idxCh)
	}
	return nil
}

func (c *consumer) PushBack(ctx context.Context, job *Job) error {
	select {
	case c.queue[job.CurrentUserID%userIdIdx] <- job: // 找到当前用户的列表进行排队
		log.Printf("[%s]push back job: uid: %d, amount: %v", job.operation, job.CurrentUserID, job.Amount)
	case <-ctx.Done():
		job.Finish()
		return errors.New("failed to consume job: " + ctx.Err().Error())
	default:
		log.Printf("[err] consumer queue is full, uid: %d", job.CurrentUserID)
	}
	return nil
}

func (c *consumer) Produce(pd *kfkmodule.PushData) {
	c.pusher.PushWallet(pd)
}
