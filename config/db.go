package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

func ConnectDb() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		GetConfigEnv("db.user"), GetConfigEnv("db.password"), GetConfigEnv("db.host"), GetConfigEnv("db.port"), GetConfigEnv("db.name"))

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db, err
}

func GetDb() *gorm.DB {
	return db
}
