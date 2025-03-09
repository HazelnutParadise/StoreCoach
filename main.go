package main

import (
	"StoreCoach/database"
	"StoreCoach/routes"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed frontend/dist/index.html
var HTMLfile []byte

//go:embed frontend/dist/assets/*
var AssetsDir embed.FS

const DevMode = false

func init() {
	err := godotenv.Load("StoreCoach.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	log.Println(os.Getenv("GEMINI_API_KEY"))
}

func main() {
	database.InitMongo(DevMode)
	subFS, err := fs.Sub(AssetsDir, "frontend/dist/assets")
	if err != nil {
		log.Fatal(err)
	}
	if !DevMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	routes.SetRoutes(r, HTMLfile, http.FS(subFS))

	r.Run(":8080")

}
