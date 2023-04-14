package main

import (
	"final-project/database"
	"final-project/router"
	"fmt"
)

func main() {
	database.StartDB()
	fmt.Println("Starting database...")
	router.New().Run(":3000")

	fmt.Println("Starting router on port 3000")
}
