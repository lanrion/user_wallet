package wallet

import (
	"encoding/json"
	"user_wallet/pkg/internal/kfkmodule"
)

type Pusher interface {
	PushWallet(d *kfkmodule.PushData)
}

type pusher struct {
	noticeChan  chan *kfkmodule.PushData // todo 分片
	kfkProducer kfkmodule.KfkProducer
}

func NewPusher(userShard int32, kfkProducer kfkmodule.KfkProducer, size int) Pusher {
	p := pusher{noticeChan: make(chan *kfkmodule.PushData, size), kfkProducer: kfkProducer}
	go func() {
		for wallet := range p.noticeChan {
			b, _ := json.Marshal(&wallet)
			err := p.kfkProducer.SendMessage(b)
			if err != nil {
				panic(err)
			}
		}
	}()
	return &p
}

func (p *pusher) PushWallet(d *kfkmodule.PushData) {
	p.noticeChan <- d
}
