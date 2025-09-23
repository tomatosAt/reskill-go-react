package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/config"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Context struct {
	Config *config.Config
	Router *fiber.App
	log    *logrus.Entry
	Tracer *sdktrace.TracerProvider
}

type App struct {
	*Context
}

func New(cfg *config.Config) *App {
	l := logrus.New()
	l.SetLevel(cfg.App.LogLevel)
	return &App{Context: &Context{
		Config: cfg,
		log:    l.WithField("package", "app"),
	}}
}
