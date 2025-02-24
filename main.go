package main

import (
	"StoreCoach/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	routes.SetRoutes(r)
	r.Run(":8080")
}
