package kfkmodule

import (
	"github.com/IBM/sarama"
	"sync"
)

var (
	clientOnce sync.Once
	saramaCli  sarama.Client
)

func fetchKfkClient(addrs []string) sarama.Client {
	if len(addrs) == 0 {
		addrs = []string{"127.0.0.1:9092"}
	}
	clientOnce.Do(func() {
		var err error
		scf := sarama.NewConfig()
		scf.Producer.Return.Successes = true
		scf.Consumer.Return.Errors = true
		scf.Metadata.AllowAutoTopicCreation = true

		saramaCli, err = sarama.NewClient(addrs, scf)
		if err != nil {
			panic(err)
		}
	})
	return saramaCli
}

func NewKafkaProducerClient(addrs []string) (sarama.SyncProducer, error) {
	return sarama.NewSyncProducerFromClient(fetchKfkClient(addrs))
}

func NewKafkaConsumerClient(addrs []string) (sarama.Consumer, error) {
	return sarama.NewConsumerFromClient(fetchKfkClient(addrs))
}
