package rabbitmq

type ExConf struct {
	ExchangeName string
	ExchangeKind string
	RouteKey     string
	QueueName    string
	AutoDelete   bool
}

func (e ExConf) ToRow() []string {
	return []string{e.ExchangeName, e.ExchangeKind, e.RouteKey, e.QueueName}
}

type Worker interface {
	Config() *ExConf
}

type Sender interface {
	Worker
	Send(data interface{}) error
}

type Receiver interface {
	Worker
	// error 	错误信息 (丢弃原因)
	// bool 	是否重试
	Do(msg []byte) (error, bool)
}

type Logger interface {
	Save(MessageLog) error
}

type MessageLog struct {
	Route   string // route key
	Message string // message body
	Result  string // err or result
	Retry   int64  // retry time
}
