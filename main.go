package main

import (
	"StoreCoach/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("GEMINI_API_KEY") == "" {
		err := godotenv.Load("StoreCoach.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	routes.SetRoutes(r)
	r.Run(":8080")

}
