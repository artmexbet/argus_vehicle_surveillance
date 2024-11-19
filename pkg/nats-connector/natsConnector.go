package nats_connector

import (
	"github.com/nats-io/nats.go"
	"time"
)

type Config struct {
	Host string `yaml:"host" env-prefix:"HOST" env-default:"localhost"`
	Port string `yaml:"port" env-prefix:"PORT" env-default:"4222"`
}

type Nats struct {
	cfg  *Config
	conn *nats.Conn
}

func New(cfg *Config) (*Nats, error) {
	n := &Nats{cfg: cfg}
	conn, err := nats.Connect(n.cfg.Host + ":" + n.cfg.Port)
	if err != nil {
		return nil, err
	}
	n.conn = conn
	return n, nil
}

func (n *Nats) CreateReader(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return n.conn.Subscribe(subject, handler)
}

func (n *Nats) Publish(subject string, data []byte) error {
	return n.conn.Publish(subject, data)
}

func (n *Nats) Request(subject string, data []byte) ([]byte, error) {
	msg, err := n.conn.Request(subject, data, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return msg.Data, err
}

func (n *Nats) Connection() *nats.Conn {
	return n.conn
}

func (n *Nats) Close() {
	n.conn.Close()
}
