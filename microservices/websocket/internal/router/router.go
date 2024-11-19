package router

import (
	"Argus/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type IService interface {
	GetChannel() chan models.YOLOJson
}

type Config struct {
	Port string `yaml:"port" env:"PORT" env-default:":8080"`
}

type Router struct {
	config *Config
	app    *fiber.App
	svc    IService
}

func New(cfg *Config, svc IService) *Router {
	app := fiber.New()
	r := &Router{
		config: cfg,
		app:    app,
		svc:    svc,
	}

	app.Get("/", websocket.New(r.GetInfo))
	return r
}

func (r *Router) Start() error {
	return r.app.Listen(r.config.Port)
}
