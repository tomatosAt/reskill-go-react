package module

import (
	"net/http"

	"github.com/tomatosAt/reskill-go-react/app"
	"github.com/tomatosAt/reskill-go-react/middleware"
	frontend "github.com/tomatosAt/reskill-go-react/module/front-end"
	openapi "github.com/tomatosAt/reskill-go-react/module/open-api"
)

func Create(app *app.Context) error {
	l := app.NewLogger().WithField("module", "generic")
	c, err := app.NewCacheClient(l)
	if err != nil {
		l.Errorln("[x] Init global caching module error -:", err)
		return err
	}

	aclSkipper := middleware.NewSkipperPath("")
	aclSkipper.Add("/api/-/health", http.MethodGet)
	app.Router.Use(middleware.NewACLMiddleware(&aclSkipper, c))

	// app.Router.Use(c)
	if err := openapi.Create(app); err != nil {
		l.Errorln("[x] Create OpenAPI module error -:", err)
		return err
	}
	// frontend web
	if err := frontend.Create(app); err != nil {
		l.Errorln("[x] Create FrontEndAPI module error -:", err)
		return err
	}
	return nil
}
