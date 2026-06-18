package main

import (
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
	// Development flags
	seedFlag := flag.Bool("seed", false, "Populates the Database with the seeds.")
	resetFlag := flag.Bool("resetDB", false, "Resets all DB data")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: .env not found. Using environmental variables.")
	}

	env := os.Getenv("APP_ENV")

	if env == "production" {
		if *resetFlag || *seedFlag {
			log.Fatalf("❌ Security Error: It's not allowed to use these flags in production.")
		}
	}

	// Database connection
	if !database.ConnectDB(*resetFlag) {
		log.Fatalln("❌ Fatal error: Could not connect to database or run migrations.")
	}

	defer database.DB.Close()

	// Database flags functions
	if *seedFlag {
		database.RunSeeds()
		return
	}

	if *resetFlag {
		return
	}

	// Port config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	// CORS config
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

	// Routes config
	routes.ConfigRoutes(router, database.DB)

	if string(port[0]) != ":" {
		port = ":" + port
	}
	router.Run(port)
}
