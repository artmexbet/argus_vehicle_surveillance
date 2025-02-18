package router

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/nats-io/nats.go"
)

type IBroker interface {
	CreateReader(string, nats.MsgHandler) (*nats.Subscription, error)
	Publish(string, []byte) error
	Request(string, []byte) ([]byte, error)
}

type Config struct {
	Host string `yaml:"host" env-prefix:"HOST" env-default:"localhost"`
	Port string `yaml:"port" env-prefix:"PORT" env-default:"8080"`
}

type Router struct {
	cfg       *Config
	app       *fiber.App
	validator *validator.Validate
	broker    IBroker
}

func New(cfg *Config, broker IBroker) *Router {
	app := fiber.New(fiber.Config{})
	router := &Router{
		cfg:       cfg,
		app:       app,
		validator: validator.New(),
		broker:    broker,
	}

	app.Use(logger.New())
	app.Use(recover.New())

	router.app.Get("/swagger/*", swagger.HandlerDefault)

	router.app.Post("/alarm", router.AlarmOn())
	router.app.Get("/cars", router.GetCars)
	return router
}

func (r *Router) Run() error {
	return r.app.Listen(fmt.Sprintf("%s:%s", r.cfg.Host, r.cfg.Port))
}
