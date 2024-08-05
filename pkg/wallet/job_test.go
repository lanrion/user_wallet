package wallet

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJob(t *testing.T) {

	t.Run("Deposit success", func(t *testing.T) {
		job := NewJob(1, 1.9, Deposit, 0)
		ctr := gomock.NewController(t)
		cusm := NewMockConsumer(ctr)
		cusm.EXPECT().Produce(gomock.Any()).Return().AnyTimes()
		job.Process(cusm)
		assert.NoError(t, job.Err)
	})

	t.Run("Withdraw failed", func(t *testing.T) {
		job := NewJob(10, 1.9, Withdraw, 0)
		ctr := gomock.NewController(t)
		cusm := NewMockConsumer(ctr)
		cusm.EXPECT().Produce(gomock.Any()).Return().AnyTimes()
		job.Process(cusm)
		assert.Error(t, job.Err)
	})

	t.Run("Withdraw success", func(t *testing.T) {
		wall := BucketManager.Fetch(10)
		err := wall.Deposit(100)
		assert.NoError(t, err)
		job := NewJob(10, 1.9, Withdraw, 0)
		ctr := gomock.NewController(t)
		cusm := NewMockConsumer(ctr)
		cusm.EXPECT().Produce(gomock.Any()).Return().AnyTimes()
		job.Process(cusm)
		job.Wait()
		assert.NoError(t, job.Err)
	})

	t.Run("Transfer failed", func(t *testing.T) {
		job := NewJob(110, 111.9, Transfer, 111)
		ctr := gomock.NewController(t)
		cusm := NewMockConsumer(ctr)
		cusm.EXPECT().Produce(gomock.Any()).Return().AnyTimes()
		job.Process(cusm)
		assert.Error(t, job.Err)
		assert.Equal(t, "not enough balance", job.Err.Error())
	})

	t.Run("Transfer success", func(t *testing.T) {
		wall := BucketManager.Fetch(210)
		_ = wall.Deposit(2000)
		job := NewJob(210, 100, Transfer, 211)
		ctr := gomock.NewController(t)
		cusm := NewMockConsumer(ctr)
		cusm.EXPECT().Produce(gomock.Any()).Return().AnyTimes()
		job.Process(cusm)
		toWall := BucketManager.Fetch(211)
		assert.Equal(t, float64(2000-100), wall.Balance())
		assert.Equal(t, float64(100), toWall.Balance())

		assert.NoError(t, job.Err)
	})

}
