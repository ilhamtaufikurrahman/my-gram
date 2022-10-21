package database

import (
	"fmt"
	"my-gram/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "123456"
	dbPort   = "5432"
	dbName   = "my_gram"
	db       *gorm.DB
	err      error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", host, user, password, dbPort, dbName)

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		fmt.Println("Error connecting database:", err)
		return
	}

	fmt.Print("Successfuly connecting to database")
	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}

func GetDB() *gorm.DB {
	return db
}
