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

// @title Argus API docs
// @version 1.0
// @description This is a sample swagger docs for argus project
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
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

	r := router.New(cfg.RouterConfig, broker)
	log.Fatal(r.Run())
}
