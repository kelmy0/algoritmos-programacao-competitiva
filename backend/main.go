package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/config"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/database"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/middleware"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/routes"
)

func main() {
	// Development flags
	seedFlag := flag.Bool("seed", false, "Populates the Database with the seeds.")
	resetFlag := flag.Bool("resetDB", false, "Resets all DB data")
	flag.Parse()

	cfg := config.LoadConfig()

	if cfg.AppEnv == "production" {
		if *resetFlag || *seedFlag {
			log.Fatalf("❌ Security Error: It's not allowed to use these flags in production.")
		}
	}
	//Google Oauth config
	googleCfg := config.LoadGoogleOauthConfig(cfg.GoogleClientId, cfg.GoogleClientSecret, "http://localhost:8000/api/auth/google/callback")

	// Database connection
	if !database.ConnectDB(*resetFlag, cfg.DatabaseURL) {
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
	router := gin.Default()

	// CORS config
	router.Use(middleware.SetupCORS(cfg))

	routes.ConfigRoutes(router, database.DB, cfg, googleCfg)
	router.Run(cfg.Port)
}
