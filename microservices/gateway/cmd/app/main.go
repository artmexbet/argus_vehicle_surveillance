package main

import (
	natsconnector "Argus/pkg/nats-connector"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	NatsConfig *natsconnector.Config `yaml:"nats" env-prefix:"NATS_"`
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

	_, err = natsconnector.New(cfg.NatsConfig)
	if err != nil {
		log.Fatal(err)
	}

}
