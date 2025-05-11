package db

import (
	"final/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("expressions.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(&models.User{}, &models.Expression{}, &models.Task{})
	if err != nil {
		log.Fatal("migration failed")
	}

	log.Println("Successfully initialized db.")
}

func GetDB() *gorm.DB {
	return db
}

func CloseDatabase() {
	if sqlDB, err := db.DB(); err != nil {
		log.Fatalf("[FATAL] [PSQL] %s", err.Error())
	} else {
		sqlDB.Close()
	}
}
