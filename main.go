package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	requestLimit  = 5
	windowSeconds = 60
)

var userRequests = make(map[string][]time.Time)
var mu sync.Mutex

// RateLimitMiddleware applies rate limiting per user
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		if isRateLimited(userID) {
			log.Printf("Rate limit exceeded for user: %s", userID)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
			return
		}

		c.Next()
	}
}

// isRateLimited checks if a user has exceeded the rate limit and cleans up old timestamps
func isRateLimited(userID string) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	reqs, exists := userRequests[userID]
	if !exists {
		reqs = []time.Time{}
	}

	// Keep only requests within the time window
	var recentRequests []time.Time
	for _, reqTime := range reqs {
		if now.Sub(reqTime) <= windowSeconds*time.Second {
			recentRequests = append(recentRequests, reqTime)
		}
	}

	recentRequests = append(recentRequests, now)
	userRequests[userID] = recentRequests

	return len(recentRequests) > requestLimit
}

func main() {
	router := gin.Default()

	// Add and configure the CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	api.Use(RateLimitMiddleware())
	{
		api.GET("/resource", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Access granted to resource!"})
		})
		api.POST("/data", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Data processed!"})
		})
	}

	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a public resource, not rate-limited."})
	})

	log.Printf("Rate Limiter Service starting on :8081...")
	log.Fatal(router.Run(":8081"))
}
