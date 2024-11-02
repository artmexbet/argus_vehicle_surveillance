package main

import (
	postgresConnector "Argus/microservices/management-service/internal/postgres-connector"
	natsConnector "Argus/pkg/nats-connector"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	NatsConfig *natsConnector.Config     `yaml:"nats" env-prefix:"NATS_"`
	DbConfig   *postgresConnector.Config `yaml:"db" env-prefix:"DB_"`
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

	broker, err := natsConnector.New(cfg.NatsConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer broker.Close()

	conn, err := postgresConnector.New(cfg.DbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
}
