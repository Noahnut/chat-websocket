package streaming

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
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
	uid            string
	mutex          sync.RWMutex
	natsAddr       string
	conn           *nats.Conn
	jetStream      jetstream.JetStream
	subjectHandler map[string]subscribeCallback
	consumer       jetstream.Consumer
	consumerConfig jetstream.ConsumerConfig
}

func NewNATS(ctx context.Context, uid, natsAddr string) *NATS {
	return &NATS{
		ctx:            ctx,
		natsAddr:       natsAddr,
		uid:            uid,
		subjectHandler: make(map[string]subscribeCallback),
	}
}

func (n *NATS) Connect() error {
	conn, err := nats.Connect(n.natsAddr)

	if err != nil {
		return err
	}

	jetStream, err := jetstream.New(conn)

	if err != nil {
		return err
	}

	n.conn = conn
	n.jetStream = jetStream

	// TODO: move to jetStream Init and config
	message_store_config := jetstream.StreamConfig{
		Name:      MESSAGE_STREAM_STORE_NAME,
		Subjects:  []string{MESSAGE_STREAM_STORE_SUBJECT},
		Retention: jetstream.WorkQueuePolicy,
		Discard:   jetstream.DiscardOld,
		Replicas:  3,
	}

	if _, err := n.jetStream.CreateOrUpdateStream(n.ctx, message_store_config); err != nil {
		return err
	}

	message_config := jetstream.StreamConfig{
		Name:      MESSAGE_STREAM_NAME,
		Subjects:  []string{MESSAGE_STREAM_SUBJECT_WILDCARD},
		Retention: jetstream.LimitsPolicy,
		Discard:   jetstream.DiscardOld,
		MaxAge:    60 * time.Minute,
		Replicas:  3,
	}

	if _, err := n.jetStream.CreateOrUpdateStream(n.ctx, message_config); err != nil {
		return err
	}

	consumerConfig := jetstream.ConsumerConfig{
		AckPolicy:      jetstream.AckExplicitPolicy,
		Durable:        n.uid,
		Name:           n.uid,
		FilterSubjects: []string{MESSAGE_STREAM_SUBJECT_TEMP + n.uid},
	}

	n.consumerConfig = consumerConfig

	consumer, err := n.jetStream.CreateOrUpdateConsumer(n.ctx, MESSAGE_STREAM_NAME, consumerConfig)

	if err != nil {
		return err
	}

	n.consumer = consumer

	go n.subscribeRoutine()

	return err
}

func (n *NATS) subscribeRoutine() {

	for {
		select {
		case <-n.ctx.Done():
			return
		default:
			msg, err := n.consumer.Fetch(10)

			if err != nil {
				log.Println(err)
				continue
			}

			for m := range msg.Messages() {
				n.mutex.RLock()
				callback, ok := n.subjectHandler[m.Subject()]

				n.mutex.RUnlock()

				m.Ack()

				if !ok {
					continue
				}

				callback(m.Data())

			}

		}
	}
}

func (n *NATS) Close() {
	n.conn.Close()
}

func (n *NATS) Publish(subject string, data []byte) error {
	if _, err := n.jetStream.Publish(n.ctx, subject, data); err != nil {
		return err
	}

	return nil
}

func (n *NATS) Subscribe(subject string, callback subscribeCallback) error {

	n.mutex.Lock()
	n.subjectHandler[subject] = callback
	n.mutex.Unlock()

	n.consumerConfig.FilterSubjects = append(n.consumerConfig.FilterSubjects, subject)

	if _, err := n.jetStream.CreateOrUpdateConsumer(n.ctx, MESSAGE_STREAM_NAME, n.consumerConfig); err != nil {
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
