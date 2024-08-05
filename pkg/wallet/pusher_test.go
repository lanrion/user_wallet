package wallet

import (
	"github.com/golang/mock/gomock"
	"testing"
	"time"
	"user_wallet/pkg/internal/kfkmodule"
)

func TestNewPusher(t *testing.T) {
	gtr := gomock.NewController(t)
	kfkProduerMock := kfkmodule.NewMockKfkProducer(gtr)

	kfkProduerMock.EXPECT().SendMessage(gomock.Any()).AnyTimes()
	pr := NewPusher(0, kfkProduerMock, 100)
	time.Sleep(time.Second)
	pr.PushWallet(&kfkmodule.PushData{})
}
