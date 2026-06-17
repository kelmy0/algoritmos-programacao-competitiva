package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/database"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/routes"
)

func main() {
	seedFlag := flag.Bool("seed", false, "Populates the Database with the seeds.")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: .env not found. Using environmental variables.")
	}

	if !database.ConnectDB() {
		log.Fatalln("❌ Fatal error: Could not connect to database or run migrations.")
	}

	defer database.DB.Close(context.Background())

	if *seedFlag {
		database.RunSeeds()
		return
	}

	env := os.Getenv("APP_ENV")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	switch env {
	case "production":
		gin.SetMode(gin.ReleaseMode)

		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"https://algoritimos-programacao-competitiva.com"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

	case "development":
		gin.SetMode(gin.DebugMode)

		router.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: false,
		}))

	default:
		log.Fatalln("❌ Fatal error: Missing .env parameters.")
	}

	routes.ConfigRoutes(router)

	if string(port[0]) != ":" {
		port = ":" + port
	}
	router.Run(port)

}
