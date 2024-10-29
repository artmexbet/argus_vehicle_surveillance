package nats_connector

import "github.com/nats-io/nats.go"

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

func (n *Nats) Connection() *nats.Conn {
	return n.conn
}
