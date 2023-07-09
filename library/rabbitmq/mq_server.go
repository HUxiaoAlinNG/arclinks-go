package rabbitmq

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/assembla/cony"
	"github.com/olekukonko/tablewriter"
	"github.com/streadway/amqp"

	"arclinks-go/library/log"
)

var (
	Server    *MQServer
	retryPbl  *cony.Publisher
	failedPbl *cony.Publisher
)

type Message struct {
	routingKey string
	exchange   string
	headers    amqp.Table
	body       []byte
	ack        func(bool) error
	nack       func(bool, bool) error
}

type MQServer struct {
	Client  *cony.Client
	workers []Worker
	logger  Logger
	backOff []int
}

type Config struct {
	Host     string
	User     string
	Password string
	Port     string
	Vhost    string
	ApiPort  string
	BackOff  string
}

func NewServer(conf Config) error {
	amqpURI := "amqp://guest:guest@localhost:5672"
	if conf.Host != "" {
		amqpURI = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Vhost)
	}
	fmt.Printf("amqpURI: %s\n", amqpURI)

	// 检查连接是否成功
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.ErrLog.Errorf("amqp.Dial error: %v\n", err)
		return err
	}
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		log.ErrLog.Errorf("channel.Channel error: %v\n", err)
		return err
	}
	defer channel.Close()

	// 实例化客户端
	client := cony.NewClient(
		cony.URL(amqpURI),
		// 自动重连 间隔时间以斐波那契数列递增
		cony.Backoff(cony.DefaultBackoff),
	)
	Server = &MQServer{
		Client:  client,
		backOff: []int{10, 100, 200},
	}

	// 定义重试策略
	if conf.BackOff != "" {
		backOff := make([]int, 0)
		backOffS := strings.Split(conf.BackOff, ",")
		for _, bs := range backOffS {
			bi, _ := strconv.Atoi(bs)
			backOff = append(backOff, bi)
		}
		Server.backOff = backOff
	}

	// 声明重试,失败 的交换器和队列
	retryPbl = DeclareRetry()
	failedPbl = DeclareFailed()

	return nil
}

func (mq *MQServer) RegisterWorkers(workers []Worker) {
	mq.workers = workers
}

func (mq *MQServer) RegisterLogger(logger Logger) {
	mq.logger = logger
}

// 运行消费端
func (mq *MQServer) RunReceiver() {
	// 结束信号 终止所有的Go routine
	done := make(chan struct{})

	for _, t := range mq.workers {
		// 注册 队列,交换器 并绑定
		exConf := t.Config()
		if receiver, ok := t.(Receiver); ok {
			cns := DeclareConsumer(exConf)
			go mq.runConsumer(cns, receiver, done)
		}
	}

	printWorker()
	<-done
}

// 运行生产端
func (mq *MQServer) RunSender() {
	cli := mq.Client
	// 注册Publisher
	for _, t := range mq.workers {
		exConf := t.Config()
		if _, ok := t.(Sender); ok {
			taskName := reflect.TypeOf(t).Name()
			pub := DeclarePublisher(taskName, exConf)
			cli.Publish(pub)
		}
	}
	// 启动Client 自动重连
	go func() {
		for cli.Loop() {
			for err := range cli.Errors() {
				fmt.Println(err)
			}
		}
	}()
}

// 启动Worker (一个消费端 对应一个或多个并发worker)
func (mq *MQServer) runConsumer(cns *cony.Consumer, r Receiver, done chan struct{}) {
	msgCh := make(chan Message, 1)

	go func() {
		for Server.Client.Loop() {
			select {
			case msg := <-cns.Deliveries():
				log.BusinessLog.Infof(
					"Received message, route key: %s, body: %s\n, consumer tag: %s",
					msg.RoutingKey,
					msg.Body,
					msg.ConsumerTag,
				)
				fmt.Printf(
					"\nReceived message, route key: %s, body: %s\n, consumer tag: %s",
					r.Config().RouteKey,
					msg.Body,
					msg.ConsumerTag,
				)
				msgCh <- Message{
					routingKey: msg.RoutingKey,
					headers:    msg.Headers,
					body:       msg.Body,
					exchange:   msg.Exchange,
					ack:        msg.Ack,
					nack:       msg.Nack,
				}

			case err := <-cns.Errors():
				log.ErrLog.Errorf("Consumer error: %v\n", err)
			case err := <-Server.Client.Errors():
				log.ErrLog.Errorf("Client error: %v\n", err)
			case <-done:
				return
			}
		}
	}()

	// TODO 多个消费者 并行处理
	go func() {
		for {
			select {
			case msg := <-msgCh:
				err, retry := r.Do(msg.body)
				retryTime := 0
				if xDeaths, ok := msg.headers["x-death"]; ok {
					xDeathsList := xDeaths.([]interface{})
					lastOne := xDeathsList[0].(amqp.Table)
					retryTime = int(lastOne["count"].(int64))
				}
				// 记录日志
				mq.logSave(msg, err, retryTime)

				// 若消息未成功 则尝试重试
				if retry {
					mq.handleErr(msg, err, retryTime)
				}

				// 确认消费
				err = msg.ack(false)
				if err != nil {
					log.ErrLog.Errorf("message Ack fail: %v\n", err)
				}
			case <-done:
				return
			}
		}
	}()
}

func (mq *MQServer) handleErr(msg Message, err error, retryTime int) {
	if err == nil {
		return
	}
	if retryTime < len(mq.backOff) {
		headers := make(amqp.Table)
		if msg.headers != nil {
			headers = msg.headers
		}
		headers["x-orig-routing-key"] = msg.routingKey
		err := retryPbl.PublishWithRoutingKey(amqp.Publishing{
			Body:       msg.body,
			Headers:    headers,
			Expiration: strconv.Itoa(mq.backOff[retryTime] * 1000),
		}, msg.routingKey)
		if err != nil {
			log.ErrLog.Errorf("retryPbl.Publish error: %v\n", err)
		}
	} else {
		err := failedPbl.Publish(amqp.Publishing{
			Body: msg.body,
		})
		if err != nil {
			log.ErrLog.Errorf("retryPbl.Publish error: %v\n", err)
		}
	}
}

func (mq *MQServer) logSave(msg Message, err error, retryTime int) {
	if err != nil {
		log.ErrLog.Errorf(
			"Consume message error: %v, route key: %s, message body: %s\n ",
			err,
			msg.routingKey,
			msg.body,
		)
	}

	if mq.logger != nil {
		logMessage := MessageLog{
			Route:   msg.routingKey,
			Message: string(msg.body),
			Retry:   int64(retryTime),
		}
		if err != nil {
			logMessage.Result = err.Error()
		}
		err = mq.logger.Save(logMessage)
		if err != nil {
			log.ErrLog.Errorf("mq.logger.save fail")
		}
	}
}

func printWorker() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ExchangeName", "ExchangeKind", "RouteKey", "QueueName"})
	for _, w := range Server.workers {
		table.Append(w.Config().ToRow())
	}
	table.Render() // Send output
}
