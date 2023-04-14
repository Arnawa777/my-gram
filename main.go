package main

import (
	"final-project/database"
	"final-project/router"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting database...")

	dbConf := database.Database{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
	}

	port := "3000"

	if dbConf.Port != "" {
		port = os.Getenv("DB_PORT")
	}

	database.StartDB(&dbConf)
	// router.New().Run(":3000")
	router.New().Run(fmt.Sprintf(":%s", port))

	fmt.Println("Starting router on port 3000")
}
