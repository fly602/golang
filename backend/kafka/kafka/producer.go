package kafka

import (
	"go-community/backend/kafka/pool"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

var acks = sarama.WaitForAll

var compression = sarama.CompressionSnappy

type Producer struct {
	Config   *sarama.Config
	Producer sarama.SyncProducer
	MsgPool  *pool.MsgPool
}

func (p *Producer) ClientConfig(partitioner sarama.PartitionerConstructor) *sarama.Config {
	config := sarama.NewConfig()

	// 设置ack等级：0，1，-1(all)
	config.Producer.RequiredAcks = acks

	// 设置分区策略：随机分配、手动分配、根据Key的hash值分配等
	config.Producer.Partitioner = partitioner

	// Producer必选项
	config.Producer.Return.Successes = true

	// 生产者数据的压缩方式
	config.Producer.Compression = compression

	// sender 拉取缓冲区策略，一般用于系统调优
	{
		// 设置缓冲区中数据到达的阈值后，sender进行拉取数据，即sender拉取一批次的数据大小
		config.Producer.Flush.Bytes = 1024
		// 设置在长时间未达到数据的阈值时，等待Fequency设置的时间后拉取数据
		config.Producer.Flush.Frequency = time.Millisecond * 5
		// 设置缓冲区缓存多少条数据后，sender拉取数据
		config.Producer.Flush.Messages = 10

	}

	// sender NetworkClient 发送消息到broker等待响应超时的重试次数
	config.Producer.Retry.Max = 3

	p.Config = config
	log.Printf("Producer  Configuration=%+v\n", config)
	return p.Config
}

func (p *Producer) NewProducer(partitioner sarama.PartitionerConstructor) sarama.SyncProducer {

	conn, err := sarama.NewSyncProducer(brokers, p.ClientConfig(partitioner))
	if err != nil {
		log.Fatal("New producer error,", err)
	}

	p.Producer = conn
	log.Println("Producer connect kafka broker success!")
	return p.Producer
}

func (p *Producer) NewMsgPool(topic string, partition string, key string) {
	p.MsgPool = &pool.MsgPool{}
	p.MsgPool.InitMsgPool(topic, partition, key, func(pm *sarama.ProducerMessage) {
		pid, offset, err := p.Producer.SendMessage(pm)
		if err != nil {
			log.Println("Send Msg err,", err)
			return
		}
		log.Printf("Send Msg Success,pid=%+v,offset=%+v\n", pid, offset)
	})
}
