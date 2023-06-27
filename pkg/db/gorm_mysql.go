package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/FTChinese/superyard/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var devConfig = &gorm.Config{
	Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	),
}

func newGormConfig(production bool) *gorm.Config {
	if production {
		return &gorm.Config{}
	}

	return devConfig
}

func NewGormDB(c config.Connect, production bool) (*gorm.DB, error) {
	dsn := buildDSN(c)
	return gorm.Open(mysql.Open(dsn), newGormConfig(production))
}

func MustNewGormDB(c config.Connect, production bool) *gorm.DB {
	db, err := NewGormDB(c, production)
	if err != nil {
		panic(err)
	}

	return db
}

func mustGormOpenExistingDB(sqlDB *sql.DB, prod bool) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), newGormConfig(prod))

	if err != nil {
		panic(err)
	}

	return db
}

type MultiGormDBs struct {
	Read   *gorm.DB
	Write  *gorm.DB
	Delete *gorm.DB
}

func MustNewMultiGormDBs(prod bool) MultiGormDBs {
	return MultiGormDBs{
		Read:   MustNewGormDB(config.MustMySQLReadConn(prod), prod),
		Write:  MustNewGormDB(config.MustMySQLWriteConn(prod), prod),
		Delete: MustNewGormDB(config.MustMySQLDeleteConn(prod), prod),
	}
}
