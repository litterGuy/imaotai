package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var gormdb *gorm.DB

func Init(dbpath string) (err error) {
	gormdb, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	return err
}
