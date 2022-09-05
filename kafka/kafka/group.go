package kafka

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type KGroup struct {
	brokers           []string
	topics            []string
	startOffset       int64
	version           string
	ready             chan bool
	group             string
	channelBuffersize int
	assignor          string
}

var assignor = sarama.StickyBalanceStrategyName
var topics = []string{
	"first",
	"zaplog",
}
var group = "ConsumerGroup-2"
var AutoCommit = false

var client sarama.Client

func NewGroup() *KGroup {
	return &KGroup{
		brokers:           brokers,
		topics:            topics,
		group:             group,
		channelBuffersize: 1000,
		ready:             make(chan bool),
		version:           "3.0.0",
		assignor:          assignor,
	}
}

func (g *KGroup) Connect() func() {
	log.Println("Group init...")
	version, err := sarama.ParseKafkaVersion(g.version)
	if err != nil {
		log.Fatal("Error Parse Kafka version,", err)
	}
	config := sarama.NewConfig()
	config.Version = version
	switch assignor {
	case sarama.StickyBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case sarama.RoundRobinBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case sarama.RangeBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}
	config.Consumer.Offsets.AutoCommit.Enable = AutoCommit
	if AutoCommit {
		config.Consumer.Offsets.AutoCommit.Interval = time.Millisecond * 20
	}

	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// fetch
	config.Consumer.Fetch.Min = 1024 * 5
	config.Consumer.Fetch.Max = 1024 * 1024 * 50

	config.Consumer.MaxWaitTime = time.Millisecond * 5

	config.ChannelBufferSize = g.channelBuffersize

	newClient, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	client = newClient
	// 获取所有topic
	topics, err := newClient.Topics()
	if err != nil {
		log.Fatal("Get Topics error,", err)
	}
	log.Println("Topics=", topics)
	for _, topic := range topics {
		pids, err := newClient.Partitions(topic)
		if err != nil {
			log.Fatal("Get Partitionid error,", topic, err)
		}
		for _, pid := range pids {
			offset, err := newClient.GetOffset(topic,
				pid,
				sarama.OffsetNewest)
			// 	time.Now().Unix()*1000-time.Second.Milliseconds()*10*60)
			if err != nil {
				log.Fatalf("Get Topics error,topic=%v, partition=%v, error=%v\n", topic, pid, err)
			}
			log.Printf("Get Offset topic=%v, partition=%v, offset=%v\n", topic, pid, offset)
		}

	}

	group, err := sarama.NewConsumerGroupFromClient(g.group, newClient)
	if err != nil {
		log.Fatal("NewConsumerGroupFromClient error,", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := group.Consume(ctx, g.topics, g); err != nil {
				// 当setup失败的时候，error会返回到这里
				log.Println("Error from consumer,", err)
				return
			}
			// check if context was cancelled, signal that the consumer should stop
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			log.Println("Group loop...")
			g.ready = make(chan bool)
		}
	}()
	<-g.ready
	log.Println("Consumer up and running...!")
	return func() {
		log.Println("Kafka Close...")
		cancel()
		wg.Wait()
		if err := group.Close(); err != nil {
			log.Println("Group close err,", err)
		}
	}
}

func (k *KGroup) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("Setup...")
	var metadata = ""
	for topic, pids := range session.Claims() {
		for _, pid := range pids {
			var offset int64
			offset, err := client.GetOffset(topic,
				pid,
				// sarama.OffsetNewest)
				time.Now().Unix()*1000-time.Second.Milliseconds()*20)
			if err != nil {
				offset = sarama.OffsetNewest
			}
			log.Printf("ResetOffSet topic=%v, partition=%v, offset=%v, metadata=%v\n", topic, pid, offset, metadata)
			session.ResetOffset(topic, pid, offset, metadata)
		}
	}
	log.Println(session.Claims())
	// Mark the consumer as ready
	close(k.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (k *KGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (k *KGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//log.Println("====>>>> to here,", session, claim)
	for message := range claim.Messages() {
		log.Printf("[topic:%s] [partiton:%d] [offset:%d] [value:%s] [time:%v]\n",
			message.Topic, message.Partition, message.Offset, string(message.Value), message.Timestamp)
		// 更新位移
		session.MarkMessage(message, "")
		if !AutoCommit {
			session.Commit()
		}
	}
	return nil
}
