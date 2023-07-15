package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Gormdb *gorm.DB

func Init(dbpath string) (err error) {
	Gormdb, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	return err
}
