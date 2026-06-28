package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/config"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/middleware"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/repositories"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

func ConfigRoutes(router *gin.Engine, db *pgxpool.Pool, cfg *config.Config) {
	isProd := cfg.AppEnv == "production"
	// Algorithm Handlers and Services
	algoRepo := repositories.NewAlgorithmRepository(db)
	algoService := services.NewAlgorithmService(algoRepo)
	algoHandler := handlers.NewAlgorithmHandler(algoService)

	//Auth
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.AppName, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	authHandler := handlers.NewAuthHandler(authService, isProd)

	api := router.Group("/api")
	{
		api.GET("/ping", handlers.AnswerPing)
		api.GET("/algorithms", algoHandler.ListAlgorithms)
		api.GET("/algorithms/:slugAndId", algoHandler.GetAlgorithm)
		api.POST("/auth", authHandler.Auth)
		api.POST("/auth/refresh", authHandler.Refresh)

		admin := api.Group("/admin")
		{
			admin.Use(middleware.Fake404Middleware(cfg.AdminHash))
			admin.Use(middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName))
			admin.Use(middleware.EmployeeMiddleware())
			admin.GET("/ping", handlers.AnswerPing)
			admin.POST("/algorithms", algoHandler.PostAlgorithm, middleware.PermissionMiddleware("create:algorithms"))
			admin.DELETE("/algorithms/:slugAndId", algoHandler.DeleteAlgorithm, middleware.PermissionMiddleware("delete:algorithms"))
			admin.PUT("/algorithms", algoHandler.PutAlgorithm, middleware.PermissionMiddleware("update:algorithms"))
		}
	}
}
