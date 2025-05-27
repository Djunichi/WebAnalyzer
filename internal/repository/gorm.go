package repository

import (
	"WebAnalyzer/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitORM(dbConf config.DBConfig) (*gorm.DB, error) {
	newLogger := gormLog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLog.Config{
			SlowThreshold:             2 * time.Second, // Slow SQL threshold
			LogLevel:                  gormLog.Silent,  // Log level
			IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,           // Don't include params in the SQL log
			Colorful:                  true,            // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dbConf.URL), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	conn, err := db.DB()
	if err != nil {
		return nil, err
	}

	// we can make following values configurable
	conn.SetConnMaxIdleTime(dbConf.MaxIdleTime)
	conn.SetConnMaxLifetime(dbConf.MaxLifetime)
	conn.SetMaxOpenConns(dbConf.MaxOpenConn)
	conn.SetMaxIdleConns(dbConf.MaxIdleConn)

	return db, nil
}
