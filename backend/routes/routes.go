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
	"golang.org/x/time/rate"
)

func ConfigRoutes(router *gin.Engine, db *pgxpool.Pool, cfg *config.Config, googleConfig, githubConfig *oauth2.Config) {
	isProd := cfg.AppEnv == "production"
	argonParams := &utils.ArgonParams{
		Memory:      cfg.Memory,
		Iterations:  cfg.Iterations,
		Parallelism: cfg.Parallelism,
		SaltLength:  cfg.SaltLength,
		KeyLength:   cfg.KeyLength,
	}
	//BODY SIZE
	oneMbSize := middleware.LimitBodySize(1 * 1024 * 1024)
	tenMbSize := middleware.LimitBodySize(10 * 1024 * 1024)

	//RATE LIMIT
	standardApiLimiter := middleware.RateLimitMiddleware(middleware.NewRateLimiter(rate.Limit(5), 10))
	authFlowLimiter := middleware.RateLimitMiddleware(middleware.NewRateLimiter(rate.Limit(0.1), 5))
	strictAbuseLimiter := middleware.RateLimitMiddleware(middleware.NewRateLimiter(rate.Limit(0.0055), 2))

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
	authSocialHandler := handlers.NewAuthSocialHandler(authService, googleConfig, githubConfig, cfg.AppDomain, cfg.FrontendUrl, isProd, cfg.JwtRefreshExpiresDays)

	//Sign up
	signUpService := services.NewSignUpService(userRepo, authRepo, *argonParams, cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.AppName, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
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
		publicStandard := api.Group("", standardApiLimiter)
		{
			publicStandard.GET("/ping", handlers.AnswerPing)
			publicStandard.GET("/algorithms", algoHandler.ListAlgorithms)
			publicStandard.GET("/algorithms/:slugAndId", algoHandler.GetAlgorithm)
		}

		auth := api.Group("/auth", oneMbSize)
		{
			authFlow := auth.Group("", authFlowLimiter)
			{
				authFlow.POST("/login", authHandler.Auth)
				authFlow.POST("/refresh", authHandler.Refresh)
				authFlow.GET("/google", authSocialHandler.GoogleLogin)
				authFlow.GET("/google/callback", authSocialHandler.GoogleCallback)
				authFlow.GET("/github", authSocialHandler.GithubLogin)
				authFlow.GET("/github/callback", authSocialHandler.GithubCallback)
				auth.POST("/sign-up", signUpHandler.SignUp)
			}

			authStrict := auth.Group("", strictAbuseLimiter)
			{
				authStrict.POST("/forgot-password", userConfigHandler.ForgotPassword)
				authStrict.POST("/reset-password", userConfigHandler.ResetPassword)
				authStrict.POST("/verify-2fa", authHandler.Verify2FA)
			}

			authenticatedAuth := auth.Group("", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName), authFlowLimiter)
			{
				authenticatedAuth.POST("/logout", authHandler.Logout)
				authenticatedAuth.POST("/logout/all", authHandler.LogoutAll)
			}
		}

		users := api.Group("/users", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName))
		{
			me := users.Group("/me")
			{
				me.GET("", authFlowLimiter, userConfigHandler.GetMyCredentials)

				password := me.Group("/password", oneMbSize, authFlowLimiter)
				{
					password.POST("/set", userConfigHandler.DefinePassword)
					password.POST("/change", userConfigHandler.ChangePassword)
				}

				twoFa := me.Group("/2fa", oneMbSize, authFlowLimiter)
				{
					twoFa.POST("/generate", twoFactorHandler.Generate2FA)
					twoFa.POST("/enable", twoFactorHandler.Enable2FA)
					twoFa.POST("/disable", twoFactorHandler.Disable2FA)
				}

				linkSocial := me.Group("/link-social", oneMbSize, authFlowLimiter)
				{
					linkSocial.GET("/google", authSocialHandler.GoogleLinkAccount)
					linkSocial.GET("/github", authSocialHandler.GithubLinkAccount)
				}

			}
		}

		admin := api.Group("/admin", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName))
		{
			admin.Use(middleware.Fake404Middleware(cfg.AdminHash))
			admin.Use(middleware.EmployeeMiddleware())
			admin.Use(tenMbSize)
			admin.Use(standardApiLimiter)

			admin.GET("/ping", handlers.AnswerPing)
			admin.POST("/algorithms", middleware.PermissionMiddleware("create:algorithms"), algoHandler.PostAlgorithm)
			admin.DELETE("/algorithms/:slugAndId", middleware.PermissionMiddleware("delete:algorithms"), algoHandler.DeleteAlgorithm)
			admin.PUT("/algorithms", middleware.PermissionMiddleware("update:algorithms"), algoHandler.PutAlgorithm)
		}
	}
}
