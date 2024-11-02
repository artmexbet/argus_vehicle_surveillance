package main

import (
	"Argus/microservices/gateway/internal/router"
	natsconnector "Argus/pkg/nats-connector"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	NatsConfig   *natsconnector.Config `yaml:"nats" env-prefix:"NATS_"`
	RouterConfig *router.Config        `yaml:"router" env-prefix:"ROUTER_"`
}

func readConfig(filename string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := readConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	broker, err := natsconnector.New(cfg.NatsConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer broker.Close()

	r := router.New(cfg.RouterConfig)
	log.Fatal(r.Run())
}
