package main

import (
	carProcessor "Argus/microservices/management-service/internal/car-processor"
	centrifugoConnector "Argus/microservices/management-service/internal/centrifugo-connector"
	"Argus/microservices/management-service/internal/handler"
	notificationsManager "Argus/microservices/management-service/internal/notifications-manager"
	postgresConnector "Argus/microservices/management-service/internal/postgres-connector"
	wsConnector "Argus/microservices/management-service/internal/ws-connector"
	natsConnector "Argus/pkg/nats-connector"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"log/slog"
)

type Config struct {
	NatsConfig         *natsConnector.Config        `yaml:"nats" env-prefix:"NATS_"`
	DbConfig           *postgresConnector.Config    `yaml:"db" env-prefix:"DB_"`
	CentrifugoConfig   *centrifugoConnector.Config  `yaml:"centrifugo" env-prefix:"CENTRIFUGO_"`
	CarProcessorConfig *carProcessor.Config         `yaml:"car-processor" env-prefix:"CAR_PROCESSOR_"`
	NotificationConfig *notificationsManager.Config `yaml:"notifications" env-prefix:"NOTIFICATIONS_"`
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
	slog.Info("Config initialised: %+v", cfg)

	broker, err := natsConnector.New(cfg.NatsConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer broker.Close()
	slog.Info("Nats broker initialised")

	db, err := postgresConnector.New(cfg.DbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	slog.Info("Postgres connector initialised")

	//ws, err := centrifugoConnector.New(cfg.CentrifugoConfig)
	//defer ws.Close()
	ws := wsConnector.New(broker) // Временно заменил centrifugo на простой сокет

	nm := notificationsManager.New(cfg.NotificationConfig, broker, db)
	cp := carProcessor.New(cfg.CarProcessorConfig, nm)

	svc := handler.New(broker, ws, db, cp)
	err = svc.Init()
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Ready. Listening the broker")

	select {}
}
