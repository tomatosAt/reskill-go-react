package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func (ctx *Context) InitFiberServer() {
	ctx.log.Infoln("[*] Initialize fiber router")
	cfg := fiber.Config{
		ReadTimeout:           ctx.Config.Server.TimeoutRead,
		WriteTimeout:          ctx.Config.Server.TimeoutWrite,
		IdleTimeout:           ctx.Config.Server.TimeoutIdle,
		DisableStartupMessage: true,
		Prefork:               false,
		ServerHeader:          ctx.Config.Server.ServerHeader,
		ProxyHeader:           ctx.Config.Server.ProxyHeader,
		ErrorHandler:          serverErrorHandler,
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
		ReadBufferSize:        ctx.Config.Server.ReadBufferSize,
		BodyLimit:             ctx.Config.Server.BodyLimit}

	r := fiber.New(cfg)
	if ctx.Config.Server.EnableCORS {
		ctx.log.Infoln("[*] Used fiber cors middleware")
		r.Use(cors.New())
	}
	ctx.log.Infoln("[*] Used fiber request id middleware")
	r.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return ksuid.New().String()
		},
	}))
	ctx.log.Infoln("[*] Used custom fiber logger middleware")
	r.Use(logger.New(logger.Config{
		Format:     "${locals:requestid} - ${ip} - ${method} ${path} ${status} - ${latency}\n",
		TimeZone:   "Asia/Bangkok",
		TimeFormat: time.ANSIC,
		Next: func(c *fiber.Ctx) bool {
			// no log for health check
			if c.Path() == "/api/-/health" {
				return true
			}
			return false
		},
	}))
	ctx.log.Infoln("[*] Use recovery on crash middleware")
	r.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: fiberStackTraceHandler,
	}))

	ctx.log.Infoln("[*] Set X-Server headers")
	r.Use(func(ctx *fiber.Ctx) error {
		h, _ := os.Hostname()
		ctx.Set("X-Server-By", strings.TrimPrefix(h, "go-health-id-"))
		userCtx := context.WithValue(ctx.UserContext(), "requestid", util.GetHttpRequestId(ctx.Context()))
		ctx.SetUserContext(userCtx)
		return ctx.Next()
	})

	if ctx.Config.Tracer.Enable {
		tp, err := ctx.InitTracer()
		if err != nil {
			ctx.log.Errorln("[x] InitTracer Error ", err)
		} else {
			ctx.Tracer = tp
		}
	}

	ctx.log.Infoln("[*] Set api router")
	ctx.Router = r

	// Init healthcheck
	ctx.log.Infoln("[*] Initialize healthcheck endpoint")
	ctx.Router.Get("/api/-/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "ok"})
	})
}

// serverErrorHandler that process return errors from handlers
var serverErrorHandler = func(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	if code >= fiber.StatusInternalServerError {
		logrus.Errorln("[PANIC] ", fmt.Sprintf("[%s]", ctx.IP()), ctx.Route().Method, ctx.Route().Path, ":", err)
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return ctx.Status(code).JSON(fiber.Map{"error": err.Error()})
}

func fiberStackTraceHandler(_ *fiber.Ctx, e interface{}) {
	logrus.Errorln(fmt.Sprintf("[PANIC] %v\n%s\n", e, debug.Stack()))
}

func (ctx *Context) StartHTTP() error {
	// print debug route
	if ctx.Config.App.IsDebug() {
		fr := ctx.Router.GetRoutes()
		for _, r := range fr {
			ctx.log.Debugln(r.Name, r.Method, r.Path)
		}
	}

	serverShutdown := make(chan os.Signal, 1)
	signal.Notify(serverShutdown, os.Interrupt)
	go func() {
		// Listen for syscall signals for process to interrupt/quit
		_ = <-serverShutdown
		ctx.log.Infoln("[*] Server terminating...")
		if ctx.Tracer != nil {
			_ = ctx.Tracer.Shutdown(context.Background())
		}
		if err := ctx.Router.Shutdown(); err != nil {
			ctx.log.Errorln(fmt.Sprintf("[x] Server shutdown failed: %+v", err))
		}
	}()

	// Run the server
	srvBound := fmt.Sprintf(
		"%s:%s",
		ctx.Config.Server.ListenIp,
		ctx.Config.Server.Port,
	)

	ctx.log.Infoln(fmt.Sprintf("[*] Starting server at %s", srvBound))
	err := ctx.Router.Listen(srvBound)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		ctx.log.Errorln("[x] Start server error:", err.Error())
		return err
	}
	return nil
}
