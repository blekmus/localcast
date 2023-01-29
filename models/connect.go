package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Tabler interface {
	TableName() string
}

func ConnectDatabase(path string) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	DB = database
}
