package gormlad

import (
	"context"
	"errors"
	"time"

	"github.com/omivix/lad"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Logger struct {
	lg                        *lad.Logger
	level                     logger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

func New(lg *lad.Logger) *Logger {
	return &Logger{
		lg:                        lg,
		level:                     logger.Warn,
		slowThreshold:             800 * time.Millisecond,
		ignoreRecordNotFoundError: true,
	}
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	cp := *l
	cp.level = level
	return &cp
}

func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level < logger.Info {
		return
	}
	l.lg.Sugar().Infow("gorm", "msg", msg, "args", args)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.level < logger.Warn {
		return
	}
	l.lg.Sugar().Warnw("gorm", "msg", msg, "args", args)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.level < logger.Error {
		return
	}
	l.lg.Sugar().Errorw("gorm", "msg", msg, "args", args)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	// error
	if err != nil && !(l.ignoreRecordNotFoundError && errors.Is(err, gorm.ErrRecordNotFound)) {
		l.lg.Error("gorm query", append(fields, zap.Error(err))...)
		return
	}

	// slow
	if l.slowThreshold > 0 && elapsed > l.slowThreshold {
		if l.level >= logger.Warn {
			l.lg.Warn("gorm slow query", fields...)
		}
		return
	}

	// normal
	if l.level >= logger.Info {
		l.lg.Info("gorm query", fields...)
	}
}
