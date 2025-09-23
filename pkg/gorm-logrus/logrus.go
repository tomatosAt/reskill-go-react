package gorm_logrus

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type Logger struct {
	Logger                *logrus.Entry
	Identifier            string
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func New(identifier string, l *logrus.Entry, slowThreshold time.Duration, debug bool) *Logger {
	return &Logger{
		Identifier:            identifier,
		Logger:                l.Dup().WithField("package", "gorm-logrus"),
		SlowThreshold:         slowThreshold,
		SkipErrRecordNotFound: !debug,
		Debug:                 debug,
	}
}

func NewWithConfig(cfg Logger) *Logger {
	return &cfg
}

func (l *Logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Infoln(l.Identifier, fmt.Sprintf(s, args))
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Warnln(l.Identifier, fmt.Sprintf(s, args))
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Errorln(l.Identifier, fmt.Sprintf(s, args))
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.Logger.WithContext(ctx).WithFields(fields).Errorf("[%s] %s [%s]", l.Identifier, sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.WithContext(ctx).WithFields(fields).Warnf("[%s] %s [%s]", l.Identifier, sql, elapsed)
		return
	}

	if l.Debug {
		l.Logger.WithContext(ctx).WithFields(fields).Debugf("[%s] %s [%s]", l.Identifier, sql, elapsed)
	}
}
