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
	authHandler := handlers.NewAuthHandler(authService, isProd, cfg.AppDomain)

	//User
	userRepo := repositories.NewUserRepository(db)

	//TwoFactor
	twoFactorService := services.NewTwoFactorService(userRepo)
	twoFactorHandler := handlers.NewTwoFactorHandler(twoFactorService)

	api := router.Group("/api")
	{
		api.GET("/ping", handlers.AnswerPing)
		api.GET("/algorithms", algoHandler.ListAlgorithms)
		api.GET("/algorithms/:slugAndId", algoHandler.GetAlgorithm)
		api.POST("/auth", authHandler.Auth)
		api.POST("/auth/refresh", authHandler.Refresh)

		// Routes that requires Auth
		api.Use(middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName))
		api.POST("/auth/logout", authHandler.Logout)
		api.POST("/auth/logout/all", authHandler.LogoutAll)

		// USERS
		users := api.Group("/users")
		{
			me := users.Group("/me")
			{
				me.POST("/2fa/generate", twoFactorHandler.Generate2FA)
				me.POST("/2fa/enable", twoFactorHandler.Enable2FA)
				me.POST("/2fa/disable", twoFactorHandler.Disable2FA)
			}
		}

		admin := api.Group("/admin")
		{
			admin.Use(middleware.Fake404Middleware(cfg.AdminHash))
			admin.Use(middleware.EmployeeMiddleware())
			admin.GET("/ping", handlers.AnswerPing)
			admin.POST("/algorithms", middleware.PermissionMiddleware("create:algorithms"), algoHandler.PostAlgorithm)
			admin.DELETE("/algorithms/:slugAndId", middleware.PermissionMiddleware("delete:algorithms"), algoHandler.DeleteAlgorithm)
			admin.PUT("/algorithms", middleware.PermissionMiddleware("update:algorithms"), algoHandler.PutAlgorithm)
		}
	}
}
