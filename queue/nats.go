package queue

import "github.com/nats-io/nats.go"

const (
	TEXT_MESSAGE_SUBJECT = "text_message"
)

type NATS struct {
	natsAddr string
	conn     *nats.Conn
}

func NewNATS(natsAddr string) *NATS {
	return &NATS{natsAddr: natsAddr}
}

func (n *NATS) Connect() error {
	conn, err := nats.Connect(n.natsAddr)

	n.conn = conn

	return err
}

func (n *NATS) Close() {
	n.conn.Close()
}

func (n *NATS) Publish(subject string, data []byte) error {
	return n.conn.Publish(subject, data)
}

func (n *NATS) GetTextMessageSubject() string {
	return TEXT_MESSAGE_SUBJECT
}
