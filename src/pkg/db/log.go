package db

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"buffersnow.com/spiritonline/pkg/log"
)

type CustomGormLogger struct {
	f log.LoggingFactory
}

func (l CustomGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l CustomGormLogger) Info(ctx context.Context, format string, args ...any) {
	l.f.Info("Command", format, args...)
}

func (l CustomGormLogger) Warn(ctx context.Context, format string, args ...any) {
	l.f.Warning("Command", format, args...)
}

func (l CustomGormLogger) Error(ctx context.Context, format string, args ...any) {
	l.f.Error("Command", format, args...)
}

func (l CustomGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsedTime := time.Since(begin)
	switch {
	case err != nil && errors.Is(err, gorm.ErrRecordNotFound):
		fcStr, _ := fc()
		l.f.Trace(log.DEBUG_DATABASE, "Command", "<Time: %.3fms> <Rows: None> %s", elapsedTime, fcStr)
	case err != nil:
		fcStr, _ := fc()
		l.f.Error("Command", "<Time: %.3fms> <Rows: None> %s", elapsedTime, err.Error(), fcStr)
	default:
		sql, rows := fc()
		l.f.Trace(log.DEBUG_DATABASE, "Command", "<Time: %.3fms> <Rows: %v> %s", elapsedTime, rows, sql)
	}
}
