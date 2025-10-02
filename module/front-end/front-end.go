package frontend

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/app"
	"github.com/tomatosAt/reskill-go-react/middleware"
	"github.com/tomatosAt/reskill-go-react/module/front-end/handler"
	"github.com/tomatosAt/reskill-go-react/module/front-end/repositories"
	"github.com/tomatosAt/reskill-go-react/module/front-end/services"
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
	prefixPath := app.Config.App.PrefixPath + "-fe"
	g := app.Router.Group(prefixPath)
	skipper := middleware.NewSkipperPath("")
	skipper.Add(prefixPath+"/v1/pre-register", http.MethodPost)

	addRouter(g, h)
	return nil
}

func addRouter(r fiber.Router, h *handler.Handler) {
	v1 := r.Group("/v1")
	preRegist := v1.Group("pre-register")
	preRegist.Post("", h.PreRegisterHandler)
}
