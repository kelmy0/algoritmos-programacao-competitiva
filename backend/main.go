package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

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
		gin.SetMode(gin.ReleaseMode)
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		slog.SetDefault(logger)

		if *resetFlag || *seedFlag {
			log.Fatalf("❌ Security Error: It's not allowed to use these flags in production.")
		}
	} else {
		gin.SetMode(gin.DebugMode)
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}

	//Google Oauth config
	googleCfg := config.LoadGoogleOauthConfig(cfg.GoogleClientId, cfg.GoogleClientSecret, cfg.GoogleCallbackUrl)
	githubCfg := config.LoadGithubOauthConfig(cfg.GithubClientId, cfg.GithubClientSecret, cfg.GithubCallbackUrl)

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

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.SetupRecovery())
	router.Use(middleware.SetupCORS(cfg.AppEnv, cfg.FrontendUrl))
	router.Use(middleware.SetupSecureHeaders())
	routes.ConfigRoutes(router, database.DB, cfg, googleCfg, githubCfg)
	router.Run(cfg.Port)
}
