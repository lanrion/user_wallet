package model

import (
	"gorm.io/driver/sqlite"
	"log"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

const (
	maxIdleConns = 32
	maxOpenConns = 32
	maxLifetime  = time.Minute
)

func Init() error {
	dbOption := &gorm.Config{}
	db1, err := gorm.Open(sqlite.Open("wallet_go.db"), dbOption)
	if err != nil {
		return err
	}
	err = db1.AutoMigrate(&User{})

	if err != nil {
		return err
	}

	db1 = db1.Debug()

	sqlDB, err := db1.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	log.Print("Initialized database")
	db = db1
	return nil
}
