package wallet

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConsumer(t *testing.T) {
	gtr := gomock.NewController(t)
	pusherMock := NewMockPusher(gtr)
	pusherMock.EXPECT().PushWallet(gomock.Any()).AnyTimes()
	csm := NewConsumer(pusherMock)
	err := csm.Start()
	assert.Nil(t, err)
}

func TestConsumer_Start(t *testing.T) {
	queue := make(map[uint32]chan *Job)
	gtr := gomock.NewController(t)
	pusherMock := NewMockPusher(gtr)
	pusherMock.EXPECT().PushWallet(gomock.Any()).AnyTimes()
	csm := &consumer{queue: queue, pusher: pusherMock}
	err := csm.Start()
	assert.Nil(t, err)
	assert.Equal(t, userIdIdx, len(csm.queue))

	job := NewJob(1, 1.9, Deposit, 0)

	err = csm.PushBack(context.TODO(), job)
	assert.Nil(t, err)

}
