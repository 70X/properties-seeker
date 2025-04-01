package storage

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

func NewLogger(file *os.File) logger.Interface {
	if file == nil {
		file = os.Stdout
	}
	gormLogger := logger.New(
		log.New(file, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	return gormLogger
}
