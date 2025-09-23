package app

import (
	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/pkg/cache"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func (ctx *Context) NewCacheClient(logger *logrus.Entry) (*cache.Redis, error) {
	logger.Infoln("[*] Initialize caching")
	c := cache.NewWithCfg(cache.Config{
		Host:     ctx.Config.Redis.Host,
		Port:     ctx.Config.Redis.Port,
		Password: ctx.Config.Redis.Password,
		Db:       util.AtoI(ctx.Config.Redis.Database, 0),
	})
	if err := c.Ping(); err != nil {
		logger.Errorln("[x] caching connection error:", err.Error())
		return nil, err
	}
	return &c, nil
}
