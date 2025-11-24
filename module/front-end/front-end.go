package frontend

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/app"
	middleware "github.com/tomatosAt/reskill-go-react/middleware"
	"github.com/tomatosAt/reskill-go-react/module/front-end/handler"
	fwMiddleware "github.com/tomatosAt/reskill-go-react/module/front-end/middleware"
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
	skipper.Add(prefixPath+"/v1/session", http.MethodPost)
	skipper.Add(prefixPath+"/v1/auth/login", http.MethodPost)
	fwMid := fwMiddleware.NewFWAuthMiddleware(&skipper, repo.DB(), repo.Cache(), repo.AppCfg().Secret.PrivateKey, repo.Cache())
	addRouter(g, h, fwMid)
	return nil
}

func addRouter(r fiber.Router, h *handler.Handler, fwMid *fwMiddleware.FrontWebMiddleware) {
	v1 := r.Group("/v1").Use(fwMid.NewAuth())
	preRegist := v1.Group("pre-register")
	preRegist.Post("", h.PreRegisterHandler)
	// 	1) POST /pre-register
	// รับข้อมูลสมัคร (first_name, last_name, email, phone)
	// insert ลง DB → status = pending
	// generate OTP → save DB/Redis + ส่ง email
	// คืน response → "message": "Pre-register successful, check email for OTP"
	// 	2) POST /session
	session := v1.Group("session")
	session.Post("", h.CreateSessionHandler)
	session.Get("", h.GetProfilesHandler)
	// session.Get("/profiles", h.GetProfilesHandler)
	// 	3) POST /login
	auth := v1.Group("auth")
	auth.Post("/login", h.LoginUserPassHandler)
}
