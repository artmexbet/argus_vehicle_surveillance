package router

import "github.com/gofiber/fiber/v2"

type Config struct {
	Host string `yaml:"host" env-prefix:"HOST" env-default:"localhost"`
	Port string `yaml:"port" env-prefix:"PORT" env-default:"8080"`
}

type Router struct {
	cfg *Config
	app *fiber.App
}

func New(cfg *Config) *Router {
	app := fiber.New(fiber.Config{})
	return &Router{cfg: cfg, app: app}
}
