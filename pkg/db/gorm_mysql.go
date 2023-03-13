package db

import (
	"github.com/FTChinese/superyard/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormDB(c config.Connect) (*gorm.DB, error) {
	dsn := buildDSN(c)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func MustNewGormDB(c config.Connect) *gorm.DB {
	db, err := NewGormDB(c)
	if err != nil {
		panic(err)
	}

	return db
}
