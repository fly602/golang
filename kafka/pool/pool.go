package pool

import (
	"log"
	"strconv"
	"sync"

	"github.com/Shopify/sarama"
)

// 创建协程处理消息
var MaxMsgGoroutine = 128

// msg channel 缓冲大小
var MaxMsgChannel = 1024 * 1024

type MsgPool struct {
	pool    *sync.WaitGroup
	msgChan chan string
}

func (p *MsgPool) InitMsgPool(topic string, partition string, key string, sendmsg func(*sarama.ProducerMessage)) {
	p.pool = &sync.WaitGroup{}
	p.msgChan = make(chan string, MaxMsgChannel)
	for i := 0; i < MaxMsgGoroutine; i++ {
		go func() {
			msg := &sarama.ProducerMessage{}
			// 设置发送信息
			if partition != "" {
				p, _ := strconv.ParseInt(partition, 10, 32)
				msg.Partition = int32(p)
			}
			if key != "" {
				msg.Key = sarama.StringEncoder(key)
			}
			msg.Topic = topic
			log.Printf("Start MsgPool Goroutines,topic=%v, partition=%v,key=%v\n", topic, partition, key)
			for message := range p.msgChan {
				msg.Value = sarama.StringEncoder(message)
				sendmsg(msg)
			}
			log.Println("MsgPool Goroutines exit ...")
			p.pool.Done()
		}()
		p.pool.Add(1)
	}
}

func (p *MsgPool) CloseMsgPool() {
	close(p.msgChan)
	p.pool.Wait()
}

func (p *MsgPool) PushMsg(message string) {
	p.msgChan <- message
}
