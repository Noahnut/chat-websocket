package streaming

import (
	"github.com/nats-io/nats.go"
)

const (
	MESSAGE_STREAM_STORE_NAME = "message_store"
	MESSAGE_STREAM_NAME       = "messages"
)

const (
	MESSAGE_STREAM_STORE_SUBJECT = "message_store.*"
	MESSAGE_STREAM_SUBJECT_TEMP  = "messages."
)

type NATS struct {
	natsAddr  string
	conn      *nats.Conn
	jetStream nats.JetStreamContext
}

func NewNATS(natsAddr string) *NATS {
	return &NATS{natsAddr: natsAddr}
}

func (n *NATS) Connect() error {
	conn, err := nats.Connect(n.natsAddr)

	if err != nil {
		return err
	}

	jetStream, err := conn.JetStream()

	if err != nil {
		return err
	}

	n.conn = conn
	n.jetStream = jetStream

	message_store_config := &nats.StreamConfig{
		Name:      MESSAGE_STREAM_STORE_NAME,
		Subjects:  []string{MESSAGE_STREAM_STORE_SUBJECT},
		Retention: nats.WorkQueuePolicy,
		Discard:   nats.DiscardOld,
		Replicas:  3,
	}

	if _, err := n.jetStream.StreamInfo(MESSAGE_STREAM_STORE_NAME); err != nil {
		if _, err := n.jetStream.AddStream(message_store_config); err != nil {
			return err
		}
	} else {
		if _, err := n.jetStream.UpdateStream(message_store_config); err != nil {
			return err
		}
	}

	message_config := &nats.StreamConfig{
		Name:      MESSAGE_STREAM_NAME,
		Subjects:  []string{MESSAGE_STREAM_SUBJECT_TEMP + "*"},
		Retention: nats.InterestPolicy,
		Discard:   nats.DiscardOld,
		Replicas:  3,
	}

	if _, err := n.jetStream.StreamInfo(MESSAGE_STREAM_NAME); err != nil {
		if _, err := n.jetStream.AddStream(message_config); err != nil {
			return err
		}
	} else {
		if _, err := n.jetStream.UpdateStream(message_config); err != nil {
			return err
		}
	}

	return err
}

func (n *NATS) Close() {
	n.conn.Close()
}

func (n *NATS) Publish(subject string, data []byte) error {
	if _, err := n.jetStream.Publish(subject, data); err != nil {
		return err
	}

	return nil
}

func (n *NATS) Subscribe(subject string, callback subscribeCallback) error {
	_, err := n.jetStream.Subscribe(subject, func(msg *nats.Msg) {
		println(msg)
		callback(msg.Data)
	}, nats.AckAll())

	if err != nil {
		return err
	}

	return nil
}

// get subject send to message store service
func (n *NATS) GetMessageStoreSubject() string {
	return MESSAGE_STREAM_STORE_SUBJECT
}

func (n *NATS) GetPrivateMessageSubject(userID string) string {
	return MESSAGE_STREAM_SUBJECT_TEMP + userID
}
