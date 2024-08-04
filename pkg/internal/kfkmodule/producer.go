package kfkmodule

import (
	"fmt"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

type KfkProducer interface {
	SendMessage(payload []byte) error
}

type producer struct {
	Topic    string
	producer sarama.SyncProducer
}

func NewProducer(userSharding int32, addrs []string) (KfkProducer, error) {
	client, err := NewKafkaProducerClient(addrs)
	if err != nil {
		return nil, err
	}
	p := &producer{producer: client, Topic: GetWalletChangeTopic(fmt.Sprint(userSharding))}
	return p, nil
}

func (pd *producer) SendMessage(payload []byte) error {
	nowTimeNs := time.Now().UnixNano()
	kMsg := &sarama.ProducerMessage{
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("send_at"),
				Value: []byte(strconv.FormatInt(nowTimeNs, 10)),
			},
		},
		Topic: pd.Topic,
		Value: sarama.ByteEncoder(payload),
	}

	_, _, err := pd.producer.SendMessage(kMsg)
	return err
}
