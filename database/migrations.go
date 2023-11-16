package database

import (
	"fmt"
	"log"

	"github.com/bparsons094/go-server-base/models"
	"github.com/bparsons094/go-server-base/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateNewDB(config utils.Config) {
	log.Println("Creating New DB")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=America/Chicago", config.DBHost, config.DBUser, config.DBPassword, config.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		log.Fatalf("Error connecting to the database server: %v", err)
	}

	createDBSQL := fmt.Sprintf("CREATE DATABASE \"%s\" WITH OWNER = \"%s\";", config.DBName, config.DBUser)
	if err := db.Exec(createDBSQL).Error; err != nil {
		log.Fatalf("Error creating the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error closing the connection: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Error closing the connection: %v", err)
	}
}

func FullReset(config utils.Config) {
	log.Println("Running Full Reset")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=America/Chicago", config.DBHost, config.DBUser, config.DBPassword, config.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	dropDBSQL := fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\" WITH (FORCE);", config.DBName)
	if err := db.Exec(dropDBSQL).Error; err != nil {
		log.Fatalf("Error dropping the database: %v", err)
	}

	createDBSQL := fmt.Sprintf("CREATE DATABASE \"%s\" WITH OWNER = \"%s\";", config.DBName, config.DBUser)
	if err := db.Exec(createDBSQL).Error; err != nil {
		log.Fatalf("Error creating the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error closing the connection: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Error closing the connection: %v", err)
	}
}

var modelsTOCreateOnly = []interface{}{}

var modelsToMigrate = []interface{}{
	&models.User{},
}

func CreateAllTables(db *gorm.DB) {
	db.DisableForeignKeyConstraintWhenMigrating = true
	log.Println("Creating Tables")
	for _, model := range modelsToMigrate {
		log.Printf("Creating table for %T", model)
		if err := db.Migrator().CreateTable(model); err != nil {
			log.Printf("Error creating table for %T: %v", model, err)
		}
	}

	for _, model := range modelsTOCreateOnly {
		log.Printf("Creating table for %T", model)
		if err := db.Migrator().CreateTable(model); err != nil {
			log.Printf("Error creating table for %T: %v", model, err)
		}
	}

	db.DisableForeignKeyConstraintWhenMigrating = false
}

func RunAllAutoMigrations(db *gorm.DB) {
	log.Println("Running Auto Migrations")
	for _, model := range modelsToMigrate {
		log.Printf("Migrating %T", model)
		if err := db.Migrator().AutoMigrate(model); err != nil {
			log.Printf("Error migrating %T: %v", model, err)
		}
	}
}

func RunAllDropAddTables(db *gorm.DB) {
	log.Println("Running Drop Add Tables")
	for _, model := range modelsToMigrate {
		log.Printf("Dropping and Adding table for %T", model)
		if err := db.Migrator().DropTable(model); err != nil {
			log.Printf("Error dropping table for %T: %v", model, err)
		}
		if err := db.Migrator().CreateTable(model); err != nil {
			log.Printf("Error creating table for %T: %v", model, err)
		}
	}
}

func InstallExtensions(db *gorm.DB) {
	log.Println("Installing Extensions")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
}
