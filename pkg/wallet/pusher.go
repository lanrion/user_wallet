package wallet

import (
	"encoding/json"
	"log"
	"user_wallet/pkg/internal/kfkmodule"
)

type Pusher struct {
	noticeCah   chan *kfkmodule.PushData // todo 分片
	kfkProducer kfkmodule.KfkProducer
}

func NewPusher(userShard int32, kfkAddrs []string, size int) *Pusher {
	kfkProducer, err := kfkmodule.NewProducer(userShard, kfkAddrs)
	if err != nil {
		log.Fatal(err)
	}
	p := Pusher{noticeCah: make(chan *kfkmodule.PushData, size), kfkProducer: kfkProducer}
	go func() {
		for wallet := range p.noticeCah {
			b, _ := json.Marshal(&wallet)
			err = p.kfkProducer.SendMessage(b)
			if err != nil {
				panic(err)
			}
		}
	}()
	return &p
}

func (p *Pusher) PushWallet(d *kfkmodule.PushData) {
	p.noticeCah <- d
}
