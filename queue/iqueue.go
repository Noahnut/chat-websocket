package queue

type IQueueSubject interface {
	GetTextMessageSubject() string
}

type IQueueOperator interface {
	Connect() error
	Publish(subject string, data []byte) error
}

type IQueue interface {
	IQueueOperator
	IQueueSubject
}
