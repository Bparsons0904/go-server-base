package database

import (
	"fmt"
	"log"

	"github.com/bparsons094/go-server-base/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config utils.Config) *gorm.DB {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

	var logLevel logger.LogLevel
	if config.DBLogging == "info" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logLevel), PrepareStmt: true})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	SetDatabase(db)

	log.Println("Connected Successfully to the Database")
	return DB
}

func SetDatabase(db *gorm.DB) {
	DB = db
}

func GetDatabase() *gorm.DB {
	return DB
}
