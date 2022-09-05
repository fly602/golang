package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	// mq 连接句柄
	conn *amqp.Connection
	// 管道
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机名称
	Exchange string
	// bind key 名称
	Key string
	// 连接信息
	Mqurl string
}

// RabbitMQ连接函数
func RabbitMqUrl() (url string) {
	// RabbitMQ分配的用户名称
	var user string = "fly"
	// RabbitMQ用户的密码
	var pwd string = "123456"
	// RabbitMQ Broker 的ip地址
	var host string = "127.0.0.1"
	// RabbitMQ Broker 监听的端口
	var port string = "45672"
	// Virture Host
	var vhost string = "guest"
	url = "amqp://" + user + ":" + pwd + "@" + host + ":" + port + "/" + vhost
	return
}

func NewRabbitMq(queuename string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queuename, Exchange: exchange, Key: key, Mqurl: RabbitMqUrl()}
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// 创建简单模式下的mq实例
func NewRabbitMqSimple(queue string) *RabbitMQ {
	rabbitmq := NewRabbitMq(queue, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabb"+"itmq!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 订阅模式创建mq实例
func NewRabbitMqPubSub(exchange string) *RabbitMQ {
	rabbitmq := NewRabbitMq("", exchange, "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a  channel!")
	return rabbitmq

}

// 路由模式
// 创建rabbitmq实例
func NewRabbitMqRouting(exchange string, routingKey string) *RabbitMQ {
	// 创建rabbitmq实例
	rabbitmq := NewRabbitMq("", exchange, routingKey)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a  channel!")
	return rabbitmq

}

func (r *RabbitMQ) PublishSimple(message string) {
	// 1、申请队列，如果队列不存在则自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		true,
		//是否自动删除
		false,
		//是否具有排他性
		true,
		//是否阻塞处理
		false,
		//额外的处理
		nil,
	)
	if err != nil {
		log.Println("QueueDeclare err,", err)
	}
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		// 如果未true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) ComsumeSimple() {
	// // 1、申请队列，如果队列不存在则自动创建，存在则跳过创建
	// q, err := r.channel.QueueDeclare(
	// 	r.QueueName,
	// 	//是否持久化
	// 	true,
	// 	//是否自动删除
	// 	false,
	// 	//是否具有排他性
	// 	true,
	// 	//是否阻塞处理
	// 	false,
	// 	//额外的处理
	// 	nil,
	// )
	// if err != nil {
	// 	log.Println("QueueDeclare err,", err)
	// }
	megs, err := r.channel.Consume(
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		true,
		// 是否独有
		false,
		// 设置为true，表示不能将同一个conenction中的消费者发送的消息传递给这个connection中的消费者
		false,
		// 队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		log.Println("QueueDeclare err,", err)
	}
	forever := make(chan bool)
	go func() {
		for d := range megs {
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (r *RabbitMQ) PublishPub(message string) {
	// 1、尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an excha"+"nge")
	r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) ReceiveSub() {
	// 1、试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an excha"+"nge")
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an excha"+"nge")
	err = r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
	)

	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (r *RabbitMQ) PublishRouting(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 要改成direct
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{ContentType: "text/plain",
			Body: []byte(message),
		})
}

func (r *RabbitMQ) ReceiveRouting() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 要改成direct
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	// 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Println("订单过期,取消订单: ", string(d.Body), "取消时间：", time.Now())
		}
	}()
	fmt.Println("退出请按 CTRL+C")
	<-forever
}

func (r *RabbitMQ) PublishTopic(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	// 2.发送消息
	_ = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message)})
}

func (r *RabbitMQ) ReceivedTopic() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	q, err := r.channel.QueueDeclare(
		"fly-queue",
		true,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	// var q = amqp.Queue{
	// 	Name: "fly-queue",
	// }
	_ = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	messages, _ := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", string(d.Body))
		}
	}()

	fmt.Println("退出请按 CTRL+C")
	<-forever
}

func ConfirmOne(confirm chan amqp.Confirmation) {
	if confirmed := <-confirm; confirmed.Ack {
		log.Println("confirmed delivery with delivery tag:", confirmed.DeliveryTag)
	} else {
		log.Println("confirmed delivery of delivery tag:", confirmed.DeliveryTag)
	}

}

func (r *RabbitMQ) PublishConfirm(message string) {
	// r.channel.ExchangeDeclare(
	// 	r.Exchange,
	// 	"direct",
	// 	true,
	// 	false,
	// 	false,
	// 	true,
	// 	nil,
	// )
	r.channel.Confirm(false)
	confirm := r.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	defer ConfirmOne(confirm)
	r.channel.Publish(
		"ex-confirm111",
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (r *RabbitMQ) ReceiveConfirm() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	q, err := r.channel.QueueDeclare(
		"",
		true,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	// var q = amqp.Queue{
	// 	Name: "fly-queue",
	// }
	_ = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	messages, _ := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", string(d.Body))
		}
	}()

	fmt.Println("退出请按 CTRL+C")
	<-forever
}

func (r *RabbitMQ) PublishReturnm(message string) {
	r.channel.ExchangeDeclare(
		"ex-return",
		"direct",
		true,
		false,
		false,
		true,
		nil,
	)
	r.channel.Confirm(false)
	ack := r.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	returns := r.channel.NotifyReturn(make(chan amqp.Return, 1))
	go func(chan amqp.Confirmation, chan amqp.Return) {
		log.Println("waiting for ack...")
		for {
			select {
			case a := <-ack:
				if a.Ack {
					log.Println("recv ack...", a)
				} else {
					log.Println("recv nack...", a)
				}
			case <-returns:
				log.Println("recv returns")
			}
		}
	}(ack, returns)
	i := 0
	for {
		r.channel.Publish(
			"ex-return",
			r.Key,
			true,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
				Expiration:  fmt.Sprintf("%d", i%10*1000),
			},
		)
		i++
		time.Sleep(time.Millisecond * 10)
	}
}

func (r *RabbitMQ) ReceiveReturn() {
	err := r.channel.ExchangeDeclare(
		"ex-return",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	q, err := r.channel.QueueDeclare(
		"return-queue",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl": 1000000,
		},
	)
	r.failOnErr(err, "Failed to declare a queue")
	// var q = amqp.Queue{
	// 	Name: "fly-queue",
	// }
	_ = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	r.channel.Qos(10, 0, false)
	messages, _ := r.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)

	go func() {
		for d := range messages {
			d.Ack(false)
			log.Printf("Received a message: %s", string(d.Body))
		}
	}()

	fmt.Println("退出请按 CTRL+C")
	<-forever
}

func (r *RabbitMQ) DLXPub(message string) {
	r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (r *RabbitMQ) DXLConsume() {
	r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	q, _ := r.channel.QueueDeclare(
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) DXLSub() {
	r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	q, _ := r.channel.QueueDeclare(
		"queue-normal",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    "ex-dxl",
			"x-dead-letter-routing-key": "dxl",
			"x-message-ttl":             30000,
			//"x-max-length-bytes":        1000,
		},
	)
	r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	// messages, _ := r.channel.Consume(
	// 	q.Name,
	// 	"",
	// 	false,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// forever := make(chan bool)

	// go func() {
	// 	for d := range messages {
	// 		log.Printf("Received a message: %s", string(d.Body))
	// 		d.Reject(false)
	// 	}
	// 	time.Sleep(time.Millisecond * 100)
	// }()

	// fmt.Println("退出请按 CTRL+C")
	// <-forever
}
