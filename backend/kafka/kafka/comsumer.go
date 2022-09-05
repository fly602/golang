package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

func NewClient() {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		log.Fatal("Consumer error,", err)
	}
	partitions, err := consumer.Partitions("third")
	if err != nil {
		log.Fatal("Partitions get failed,", err)
	}
	log.Printf("Partitions list =%+v\n", partitions)
	for partition := range partitions {
		pc, err := consumer.ConsumePartition("third", int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("Failed to start consumer for partition %v,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步处理每个消息
		go func(sarama.PartitionConsumer) {
			log.Printf("PartitionConsumer = %+v\n", pc)
			for msg := range pc.Messages() {
				log.Printf("Partition:%d Offset:%v Key:%v Value:%v\n",
					msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			}

		}(pc)
	}
	select {}
}
