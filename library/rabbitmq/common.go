package rabbitmq

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/assembla/cony"
	"github.com/streadway/amqp"
)

var (
	publishers map[string]*cony.Publisher
)

func GetPublisher(pub interface{}) (*cony.Publisher, error) {
	if publishers == nil {
		publishers = make(map[string]*cony.Publisher)
	}
	var key string
	switch pub := pub.(type) {
	case Sender:
		key = reflect.TypeOf(pub).Name()
	case string:
		key = pub
	}
	publisher, ok := publishers[key]
	if !ok {
		return nil, errors.New("publisher is not register : %s\n" + reflect.TypeOf(pub).Name())
	}
	return publisher, nil
}

func Publish(pubKey interface{}, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	message := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         jsonData,
	}

	// 处理超时则返回connection timeout
	ch := make(chan error)
	pub, err := GetPublisher(pubKey)
	if err != nil {
		return err
	}
	go func() {
		err := pub.Publish(message)
		ch <- err
	}()
	timeout := time.NewTimer(5 * time.Millisecond)
	select {
	case err = <-ch:
	case <-timeout.C:
		err = errors.New("connection timeout")
	}
	timeout.Reset(0)
	close(ch)
	return err
}

// 声明发布者
func DeclarePublisher(taskName string, exConf *ExConf) *cony.Publisher {
	if publishers == nil {
		publishers = make(map[string]*cony.Publisher, 10)
	}

	exc := cony.Exchange{
		Name:    exConf.ExchangeName,
		Kind:    exConf.ExchangeKind,
		Durable: true,
	}
	Server.Client.Declare([]cony.Declaration{
		cony.DeclareExchange(exc),
	})
	pub := cony.NewPublisher(exc.Name, exConf.RouteKey)
	publishers[taskName] = pub
	return pub
}

// 声明消费者
func DeclareConsumer(exConf *ExConf) *cony.Consumer {
	exc := cony.Exchange{
		Name:       exConf.ExchangeName,
		Kind:       exConf.ExchangeKind,
		Durable:    true,
		AutoDelete: exConf.AutoDelete,
	}
	que := &cony.Queue{
		Name:       exConf.QueueName,
		Durable:    true,
		AutoDelete: exConf.AutoDelete,
	}
	bnd := cony.Binding{
		Queue:    que,
		Exchange: exc,
		Key:      exConf.RouteKey,
	}
	cli := Server.Client
	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(que),
		cony.DeclareExchange(exc),
		cony.DeclareBinding(bnd),
	})
	cns := cony.NewConsumer(que)
	cli.Consume(cns)
	return cns
}

// 声明重试交换器和队列
func DeclareRetry() *cony.Publisher {
	routeExc := cony.Exchange{
		Name:    "hook",
		Kind:    "topic",
		Durable: true,
	}
	retryExc := cony.Exchange{
		Name:    "retry",
		Kind:    "fanout",
		Durable: true,
	}
	retryQue := &cony.Queue{
		Name:    "retry_que",
		Durable: true,
		Args: amqp.Table{
			"x-dead-letter-exchange": "hook",
		},
	}
	Server.Client.Declare([]cony.Declaration{
		cony.DeclareExchange(routeExc),
		cony.DeclareExchange(retryExc),
		cony.DeclareQueue(retryQue),
		cony.DeclareBinding(cony.Binding{
			Queue:    retryQue,
			Exchange: retryExc,
		}),
	})

	// Declare and register a publisher
	// with the cony client
	pbl := cony.NewPublisher(retryExc.Name, "")
	//pbl := cony.NewPublisher(failedExc.Name, "")
	Server.Client.Publish(pbl)

	return pbl
}

// 声明重失败换器和队列
func DeclareFailed() *cony.Publisher {
	failedExc := cony.Exchange{
		Name:    "failed",
		Kind:    "fanout",
		Durable: true,
	}
	failedQue := &cony.Queue{
		Name:    "failed_queue",
		Durable: true,
		Args: amqp.Table{
			// 过期时间设置为3小时
			"x-message-ttl": 3 * 1000 * 3600,
		},
	}
	Server.Client.Declare([]cony.Declaration{
		cony.DeclareExchange(failedExc),
		cony.DeclareQueue(failedQue),
		cony.DeclareBinding(cony.Binding{
			Queue:    failedQue,
			Exchange: failedExc,
		}),
	})

	// Declare and register a publisher
	// with the cony client
	pbl := cony.NewPublisher(failedExc.Name, "")
	Server.Client.Publish(pbl)

	return pbl
}
