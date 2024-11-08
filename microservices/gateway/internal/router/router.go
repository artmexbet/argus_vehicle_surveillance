package router

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IService interface {
}

type Config struct {
	Host string `yaml:"host" env-prefix:"HOST" env-default:"localhost"`
	Port string `yaml:"port" env-prefix:"PORT" env-default:"8080"`
}

type Router struct {
	cfg       *Config
	app       *fiber.App
	validator *validator.Validate
	service   IService
}

func New(cfg *Config, service IService) *Router {
	app := fiber.New(fiber.Config{})
	router := &Router{
		cfg:       cfg,
		app:       app,
		validator: validator.New(),
		service:   service,
	}

	router.app.Get("/camera/list", router.CameraList())
	router.app.Post("/alarm", router.AlarmOn())
	return router
}

func (r *Router) Run() error {
	return r.app.Listen(fmt.Sprintf("%s:%s", r.cfg.Host, r.cfg.Port))
}
