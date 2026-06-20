package middleware

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/config"
)

func SetupCORS(cfg *config.Config) gin.HandlerFunc {
	switch cfg.AppEnv {
	case "production":
		gin.SetMode(gin.ReleaseMode)

		return cors.New(cors.Config{
			AllowOrigins:     []string{"https://algoritimos-programacao-competitiva.com"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})

	case "development":
		gin.SetMode(gin.DebugMode)

		return cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: false,
		})

	default:
		log.Fatalln("❌ Fatal error: Invalid or missing APP_ENV parameter.")
		return nil
	}
}
