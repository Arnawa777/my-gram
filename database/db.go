package database

import (
	"final-project/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Host     string
	Username string
	Password string
	Port     string
	Name     string
}

const (
	DEBUG_MODE = true
)

var (
	DB_HOST     = "localhost"
	DB_USER     = "postgres"
	DB_PASSWORD = "admin"
	DB_NAME     = "my-garm"
	DB_PORT     = 5432
	db          *gorm.DB
	err         error
)

func StartDB(conf *Database) {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

	if conf.Host != "" {
		config = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", conf.Host, conf.Username, conf.Password, conf.Name, conf.Port)
	}

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate Database
	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}

func GetDB() *gorm.DB {
	if DEBUG_MODE {
		return db.Debug()
	}

	return db
}
