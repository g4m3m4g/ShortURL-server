package shortener

import (
	"math/rand"
	"net/http"
	"simpleurl/models"
	"simpleurl/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HandleShortenRequest handles the POST request to shorten a URL
func HandleShortenRequest(c *gin.Context, db *gorm.DB) {
	var data struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize the URL to include the protocol
	originalUrl := utils.NormalizeURL(data.URL)

	// Check if the URL is accessible
	if !utils.IsURLAccessible(originalUrl) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The URL is not accessible"})
		return
	}

	var link models.SimpleUrl
	result := db.Where("original_url = ?", originalUrl).First(&link)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Generate a unique short URL
			shortUrl := generateShortUrl()
			for {
				var existingLink models.SimpleUrl
				result := db.Where("short_url = ?", shortUrl).First(&existingLink)
				if result.Error == gorm.ErrRecordNotFound {
					break
				}
				shortUrl = generateShortUrl()
			}

			link = models.SimpleUrl{OriginalUrl: originalUrl, ShortUrl: shortUrl}
			result := db.Create(&link)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"short_url": link.ShortUrl})
}

// HandleRedirect handles the redirection for a short URL
func HandleRedirect(c *gin.Context, db *gorm.DB) {
	shortUrl := c.Param("shortUrl")
	var link models.SimpleUrl

	result := db.Where("short_url = ?", shortUrl).First(&link)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, link.OriginalUrl)
}

// generateShortUrl generates a random 10-character short URL
func generateShortUrl() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	var shortUrl string
	for i := 0; i < 10; i++ {
		shortUrl += string(charset[rand.Intn(len(charset))])
	}

	return shortUrl
}
