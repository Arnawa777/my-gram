package main

import (
	"final-project/database"
	"final-project/router"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting database...")

	// port := "3000"

	port := os.Getenv("DB_PORT")
	// if dbConf.Port != "" {
	// }

	database.StartDB()
	// router.New().Run(":3000")
	router.New().Run(fmt.Sprintf(":%s", port))

	fmt.Println("Starting router on port")
}
