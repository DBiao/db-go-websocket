package kafka

import (
	"db-go-websocket/internal/global"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type ConsumerCallback func(data []byte)

func InitKafka() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	master, err := sarama.NewConsumer(global.CONFIG.Kafka.Brokers, config)
	if err != nil {
		return master, err
	}

	return master, nil
}

// ConsumerMsg 消费消息，通过回调函数进行
func ConsumerMsg(topic string, partition int32, callBack ConsumerCallback) {
	partitionConsumer, err := global.KAFKA.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if nil != err {
		global.LOG.Error("iConsumePartition error", zap.Error(err))
		return
	}

	defer partitionConsumer.Close()
	for {
		select {
		case err = <-partitionConsumer.Errors():
			global.LOG.Info("", zap.Error(err))
		case msg := <-partitionConsumer.Messages():
			if nil != callBack {
				callBack(msg.Value)
			}
		}
	}
}
