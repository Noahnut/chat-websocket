package streaming

type subscribeCallback func(msg []byte)

type IStreamSubject interface {
	GetMessageStoreSubject() string
	GetPrivateMessageSubject(userID string) string
}

type IStreamingOperator interface {
	Connect() error
	Publish(subject string, data []byte) error
	Subscribe(subject string, callback subscribeCallback) error
}

type IStreaming interface {
	IStreamingOperator
	IStreamSubject
}
