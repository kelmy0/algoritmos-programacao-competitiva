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
	authAdminHandler := handlers.NewAuthAdminHandler(authService, isProd)
	authHandler := handlers.NewAuthHandler(authService, isProd)

	api := router.Group("/api")
	{
		api.GET("/ping", handlers.AnswerPing)
		api.GET("/algorithms", algoHandler.ListAlgorithms)
		api.GET("/algorithms/:slugAndId", algoHandler.GetAlgorithm)
		api.POST("/auth", handlers.AnswerPing)
		api.POST("/auth/refresh", authHandler.Refresh)

		admin := api.Group("/admin", middleware.Fake404Middleware(cfg.AdminHash))
		{
			admin.POST("/auth", authAdminHandler.Auth)

			protectedAdmin := admin.Group("/protected", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName))
			{
				protectedAdmin.Use(middleware.EmployeeMiddleware())
				protectedAdmin.GET("/ping", handlers.AnswerPing)
				protectedAdmin.POST("/algorithms", algoHandler.PostAlgorithm, middleware.PermissionMiddleware("create:algorithms"))
				protectedAdmin.DELETE("/algorithms/:slugAndId", algoHandler.DeleteAlgorithm, middleware.PermissionMiddleware("delete:algorithms"))
				protectedAdmin.PUT("/algorithms", algoHandler.PutAlgorithm, middleware.PermissionMiddleware("update:algorithms"))
			}
		}
	}
}
