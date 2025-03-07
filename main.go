package main

import (
	"StoreCoach/routes"
	"embed"
	"io/fs"
	"log"
	"net/http"

	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed frontend/dist/index.html
var HTMLfile []byte

//go:embed frontend/dist/assets/*
var AssetsDir embed.FS

func init() {
	err := godotenv.Load("StoreCoach.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	subFS, err := fs.Sub(AssetsDir, "frontend/dist/assets")
	if err != nil {
		log.Fatal(err)
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	routes.SetRoutes(r, HTMLfile, http.FS(subFS))

	r.Run(":8080")

}
