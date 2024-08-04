package kfkmodule

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	otrace "go.opentelemetry.io/otel/trace"
	"user_wallet/pkg/internal/trace"
)

type Consumer struct {
	Topic     string
	sc        sarama.Consumer
	kfkOffset int64
}

type consumerHandler func(ctx context.Context, message *sarama.ConsumerMessage)

func NewConsumer(userSharding int32, addrs []string, kfkOfsset int64) (*Consumer, error) {
	sc, err := NewKafkaConsumerClient(addrs)
	if err != nil {
		return nil, err
	}
	cc := &Consumer{
		kfkOffset: kfkOfsset,
		sc:        sc,
		Topic:     GetWalletChangeTopic(fmt.Sprint(userSharding))}
	return cc, nil
}

func (cs *Consumer) Start(handleMsg consumerHandler) error {
	partition, err := cs.sc.ConsumePartition(cs.Topic, 0, cs.kfkOffset)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-partition.Messages():
			ctx := context.Background() // TODO 从生产者继承trace
			trace.DoWithSpan(ctx, "kafka.consumer.handler", func(ctx context.Context, span otrace.Span) {
				handleMsg(ctx, msg)
			})
		}
	}

}
