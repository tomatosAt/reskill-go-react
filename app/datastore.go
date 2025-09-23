package app

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/pkg/database"
	gorm_logrus "github.com/tomatosAt/reskill-go-react/pkg/gorm-logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
)

func (ctx *Context) NewMariaDBClient(identifier string, h string, p string, user string, pass string, name string, debug bool, logger *logrus.Entry, migration bool) (*database.Client, error) {
	logger.Infoln("[*] Initialize database", identifier)
	db := database.NewWithConfig(
		database.Config{
			Host:      h,
			Port:      p,
			Username:  user,
			Password:  pass,
			Name:      name,
			DebugMode: debug,
			Migration: migration,
		},
		ctx.NewLogger(),
	)
	logger.Infoln("[*] Connecting to database", identifier, "...")
	if err := db.ConnectWithGormConfig(gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gorm_logrus.New(identifier, logger, time.Second, ctx.Config.App.LogLevel == logrus.DebugLevel),
	}); err != nil {
		logger.Errorln("[x] database", identifier, "connection error:", err.Error())
		return nil, err
	}
	logger.Infoln("[*] Connected to database", identifier, ctx.Config.DBMain.Host, "name", ctx.Config.DBMain.Database)
	if ctx.Tracer != nil {
		_ = db.Ctx().Use(otelgorm.NewPlugin())
	}
	return &db, nil
}

func (ctx *Context) NewDBMainClient(logger *logrus.Entry) (*database.Client, error) {
	return ctx.NewMariaDBClient(
		"DBMain",
		ctx.Config.DBMain.Host,
		ctx.Config.DBMain.Port,
		ctx.Config.DBMain.User,
		ctx.Config.DBMain.Password,
		ctx.Config.DBMain.Database,
		ctx.Config.App.IsDebug(),
		logger,
		ctx.Config.DBMain.Migration,
	)
}

func (ctx *Context) NewDBLogClient(logger *logrus.Entry) (*database.Client, error) {
	return ctx.NewMariaDBClient(
		"DBLog",
		ctx.Config.DBLog.Host,
		ctx.Config.DBLog.Port,
		ctx.Config.DBLog.User,
		ctx.Config.DBLog.Password,
		ctx.Config.DBLog.Database,
		ctx.Config.App.IsDebug(),
		logger,
		ctx.Config.DBLog.Migration,
	)
}
