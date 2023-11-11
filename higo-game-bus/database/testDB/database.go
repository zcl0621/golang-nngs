package testDB

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

// db连接测试用例

var DB *gorm.DB

func GetInstance(host string, port string, user string, db_name string, password string, mode string) *gorm.DB {
	if DB == nil {
		connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, db_name, password, "disable")
		DB, _ = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
		sqlDB, err := DB.DB()
		if err != nil {
			log.Panic(err.Error())
			return nil
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxIdleTime(time.Minute * 20)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	if mode == "debug" || mode == "dev" {
		DB = DB.Debug()
	}
	return DB
}
