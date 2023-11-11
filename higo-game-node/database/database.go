package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"higo-game-node/config"
	"log"
	"time"
)

var DB *gorm.DB

func GetInstance() *gorm.DB {
	if DB == nil {
		connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", config.Conf.DataBase.Host, config.Conf.DataBase.Port, config.Conf.DataBase.User, config.Conf.DataBase.DBName, config.Conf.DataBase.Password, "disable")
		DB, _ = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
		sqlDB, err := DB.DB()
		if err != nil {
			log.Panic(err.Error())
			return nil
		}
		if config.RunMode == "debug" || config.RunMode == "dev" {
			sqlDB.SetMaxOpenConns(10)
			sqlDB.SetMaxIdleConns(5)
		} else {
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetMaxIdleConns(50)
		}
		sqlDB.SetConnMaxIdleTime(time.Second * 30)
		sqlDB.SetConnMaxLifetime(time.Second * 30)
	}
	if config.RunMode == "debug" || config.RunMode == "dev" {
		DB = DB.Debug()
	}
	return DB
}
