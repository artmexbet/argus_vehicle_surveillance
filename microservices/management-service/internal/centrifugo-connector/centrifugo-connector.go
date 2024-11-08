package centrifugo_connector

import (
	"context"
	"fmt"
	"github.com/centrifugal/centrifuge-go"
)

type Config struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"PORT" env-default:"8000"`
}

type Centrifugo struct {
	cfg    *Config
	client *centrifuge.Client
}

func New(cfg *Config) (*Centrifugo, error) {
	client := centrifuge.NewJsonClient(fmt.Sprintf("ws://%s:%s/connection/websocket", cfg.Host, cfg.Port), centrifuge.Config{})
	return &Centrifugo{
		cfg:    cfg,
		client: client,
	}, nil
}

func (c *Centrifugo) Publish(channel string, data []byte) error {
	_, err := c.client.Publish(context.Background(), channel, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Centrifugo) Close() {
	c.client.Close()
}
