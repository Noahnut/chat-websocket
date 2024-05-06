package streaming

import (
	"context"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	MESSAGE_STREAM_STORE_NAME = "message_store"
	MESSAGE_STREAM_NAME       = "messages"
)

const (
	MESSAGE_STREAM_STORE_SUBJECT    = "message_store.*"
	MESSAGE_STREAM_SUBJECT_TEMP     = "messages."
	MESSAGE_STREAM_SUBJECT_WILDCARD = "messages.*"
)

type NATS struct {
	ctx            context.Context
	mutex          sync.RWMutex
	natsAddr       string
	conn           *nats.Conn
	jetStream      nats.JetStreamContext
	subjectHandler map[string]subscribeCallback
	ch             chan *nats.Msg
}

func NewNATS(ctx context.Context, natsAddr string) *NATS {
	return &NATS{
		ctx:            ctx,
		natsAddr:       natsAddr,
		subjectHandler: make(map[string]subscribeCallback),
		ch:             make(chan *nats.Msg, 1),
	}
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

	// TODO: move to jetStream Init and config
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
		Subjects:  []string{MESSAGE_STREAM_SUBJECT_WILDCARD},
		Retention: nats.LimitsPolicy,
		Discard:   nats.DiscardOld,
		MaxAge:    60 * time.Minute,
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

	go n.subscribeRoutine()

	return err
}

func (n *NATS) subscribeRoutine() {
	for {
		select {
		case <-n.ctx.Done():
			return
		case msg := <-n.ch:
			n.mutex.RLock()
			callback, ok := n.subjectHandler[msg.Subject]
			n.mutex.RUnlock()

			if ok {
				callback(msg.Data)
				msg.Ack()
			}
		}
	}
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
	n.mutex.Lock()
	n.subjectHandler[subject] = callback
	n.mutex.Unlock()

	if _, err := n.jetStream.ChanSubscribe(subject, n.ch, nats.AckExplicit()); err != nil {
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
