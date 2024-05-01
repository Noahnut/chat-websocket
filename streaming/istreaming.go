package streaming

type IStreamSubject interface {
	GetMessageStoreSubject() string
}

type IStreamingOperator interface {
	Connect() error
	Publish(subject string, data []byte) error
}

type IStreaming interface {
	IStreamingOperator
	IStreamSubject
}
