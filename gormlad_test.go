package gormlad_test

import (
	"testing"

	"github.com/omivix/gormlad"
	"github.com/omivix/lad"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestErrorLog(t *testing.T) {
	lad.InitGlobal(
		lad.WithConsole(lad.ConsoleConfig{
			Level:   zap.DebugLevel,
			Colored: true,
		}),
		lad.WithFile(lad.FileConfig{
			Level:      zap.InfoLevel,
			Filename:   "./logs/app.log",
			MaxSizeMB:  200,
			MaxBackups: 10,
			MaxAgeDays: 30,
			Compress:   true,
			Encoding:   lad.JSONEncoding, // default is JSONEncoding
		}),
	)

	defer func() { _ = lad.Sync(lad.L()) }()

	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlad.New(lad.L()).LogMode(logger.Info),
	})
	if err != nil {
		lad.L().Error("connect mysql error", lad.Error(err))
		return
	}

	gormlad.New(lad.L())
}
