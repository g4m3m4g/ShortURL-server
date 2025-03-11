package main

import (
	//"fmt"
	"log"
	"simpleurl/config"
	"simpleurl/database"
	"simpleurl/models"
	"simpleurl/shortener"
	//"simpleurl/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
	"net/http"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	// Setup database connection
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate the model
	db.AutoMigrate(&models.SimpleUrl{})

	// Initialize the Gin router
	r := gin.Default()
	r.Use(cors.Default())

	// Setup routes
	setupRoutes(r, db)

	// Start the server
	r.Run(":8080")
}

func setupRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/shorten", func(c *gin.Context) {
		// Handle URL shortening
		shortener.HandleShortenRequest(c, db)
	})

	r.GET("/:shortUrl", func(c *gin.Context) {
		// Handle URL redirection
		shortener.HandleRedirect(c, db)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
