package kafka

import (
	"db-go-websocket/internal/global"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"log"
)

func InitKafka() error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	master, err := sarama.NewConsumer(global.CONFIG.Kafka.Brokers, config)
	if err != nil {
		return err
	}
	consumer, err := master.ConsumePartition("topic", 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	defer consumer.Close()
	go func() {
		for {
			select {
			case err = <-consumer.Errors():
				global.LOG.Info("", zap.Error(err))
			case msg := <-consumer.Messages():
				log.Println("Received messages", string(msg.Key), string(msg.Value))
			}
		}
	}()
	return nil
}
