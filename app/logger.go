package app

import (
	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func (ctx *Context) NewLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(ctx.Config.App.LogLevel)
	return l
}

func (ctx *Context) AddSyslogHook(l *logrus.Entry, tag string) *logrus.Entry {
	e := l.Dup()
	if ctx.Config.Syslog.Enable {
		if !util.Contains([]string{"udp", "tcp"}, ctx.Config.Syslog.Protocol) {
			e.Errorln("Syslog protocol error -: should be udp or tcp but got", ctx.Config.Syslog.Protocol)
			return e
		}
	}
	return e
}
