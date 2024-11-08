package main

import (
	centrifugo_connector "Argus/microservices/management-service/internal/centrifugo-connector"
	postgresConnector "Argus/microservices/management-service/internal/postgres-connector"
	"Argus/microservices/management-service/internal/service"
	natsConnector "Argus/pkg/nats-connector"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	NatsConfig       *natsConnector.Config        `yaml:"nats" env-prefix:"NATS_"`
	DbConfig         *postgresConnector.Config    `yaml:"db" env-prefix:"DB_"`
	CentrifugoConfig *centrifugo_connector.Config `yaml:"centrifugo" env-prefix:"CENTRIFUGO_"`
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

	ws, err := centrifugo_connector.New(cfg.CentrifugoConfig)
	defer ws.Close()

	svc := service.New(broker, ws)
	err = svc.Init()
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
