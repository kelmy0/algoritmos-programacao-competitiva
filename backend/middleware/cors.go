package middleware

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS(env, froendUrl string) gin.HandlerFunc {
	switch env {
	case "production":
		return cors.New(cors.Config{
			AllowOrigins:     []string{froendUrl},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})

	case "development":
		return cors.New(cors.Config{
			AllowOrigins:     []string{froendUrl},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
		})

	default:
		log.Fatalln("❌ Fatal error: Invalid or missing APP_ENV parameter.")
		return nil
	}
}
