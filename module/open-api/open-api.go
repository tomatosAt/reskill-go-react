package openapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/app"
	"github.com/tomatosAt/reskill-go-react/module/open-api/handler"
	"github.com/tomatosAt/reskill-go-react/module/open-api/repositories"
	"github.com/tomatosAt/reskill-go-react/module/open-api/services"
)

func Create(app *app.Context) error {
	repo, err := repositories.New(app)
	if err != nil {
		return err
	}
	// services
	svc := services.New(repo)
	// handler
	h := handler.NewHandler(svc)
	g := app.Router.Group(app.Config.App.PrefixPath)
	addRouter(g, h)
	return nil
}

func addRouter(r fiber.Router, h *handler.Handler) {
	v1 := r.Group("/v1")
	v1.Get("/health", h.HealthCheck)
	// v1.Post("/user", h.)
}
