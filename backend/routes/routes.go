package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/config"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/middleware"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/repositories"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"golang.org/x/oauth2"
)

func ConfigRoutes(router *gin.Engine, db *pgxpool.Pool, cfg *config.Config, googleConfig *oauth2.Config) {
	isProd := cfg.AppEnv == "production"
	argonParams := &utils.ArgonParams{
		Memory:      cfg.Memory,
		Iterations:  cfg.Iterations,
		Parallelism: cfg.Parallelism,
		SaltLength:  cfg.SaltLength,
		KeyLength:   cfg.KeyLength,
	}

	// Algorithm Handlers and Services
	algoRepo := repositories.NewAlgorithmRepository(db)
	algoService := services.NewAlgorithmService(algoRepo)
	algoHandler := handlers.NewAlgorithmHandler(algoService)

	//User
	userRepo := repositories.NewUserRepository(db)

	//Auth
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, userRepo, cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.AppName, cfg.EncryptSecretKey, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	authHandler := handlers.NewAuthHandler(authService, isProd, cfg.AppDomain, cfg.JwtRefreshExpiresDays)
	authGoogleHandler := handlers.NewAuthGoogleHandler(authService, googleConfig, cfg.AppDomain, isProd, cfg.JwtRefreshExpiresDays)

	//Sign up
	signUpService := services.NewSignUpService(userRepo, authRepo,
		cfg.Parallelism, cfg.Memory, cfg.Iterations, cfg.SaltLength,
		cfg.KeyLength, cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.AppName,
		cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	signUpHandler := handlers.NewSignUpHandler(signUpService, cfg.JwtRefreshExpiresDays, cfg.AppName, isProd)

	//TwoFactor
	twoFactorService := services.NewTwoFactorService(userRepo, cfg.EncryptSecretKey, cfg.AppName)
	twoFactorHandler := handlers.NewTwoFactorHandler(twoFactorService)

	//UserConfig
	emailService := services.NewEmailService(cfg.HostEmail, cfg.PortEmail, cfg.UserEmail, cfg.PasswordEmail, cfg.FromEmail, cfg.FrontendUrl, cfg.AppName)
	userConfigService := services.NewUserConfigService(userRepo, authRepo, *emailService, *argonParams, cfg.JwtRefreshSecret, cfg.AppName)
	userConfigHandler := handlers.NewUserConfigHandler(userConfigService)

	api := router.Group("/api")
	{
		api.GET("/ping", handlers.AnswerPing)
		api.GET("/algorithms", algoHandler.ListAlgorithms)
		api.GET("/algorithms/:slugAndId", algoHandler.GetAlgorithm)

		api.POST("/sign-up", signUpHandler.SignUp)

		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Auth)
			auth.POST("/verify-2fa", authHandler.Verify2FA)
			auth.POST("/refresh", authHandler.Refresh)
			auth.GET("/google", authGoogleHandler.GoogleLogin)
			auth.GET("/google/callback", authGoogleHandler.GoogleCallback)
			auth.POST("/forgot-password", userConfigHandler.ForgotPassword)
			auth.POST("/reset-password", userConfigHandler.ResetPassword)

			authenticatedAuth := auth.Group("", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName))
			{
				authenticatedAuth.POST("/logout", authHandler.Logout)
				authenticatedAuth.POST("/logout/all", authHandler.LogoutAll)
				authenticatedAuth.POST("/change-password", userConfigHandler.ChangePassword)
				authenticatedAuth.POST("/set-password", userConfigHandler.DefinePassword)
			}
		}

		// Routes that requires Auth

		// USERS
		users := api.Group("/users", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName))
		{
			me := users.Group("/me")
			{
				me.POST("/2fa/generate", twoFactorHandler.Generate2FA)
				me.POST("/2fa/enable", twoFactorHandler.Enable2FA)
				me.POST("/2fa/disable", twoFactorHandler.Disable2FA)
			}
		}

		admin := api.Group("/admin", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName))
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
