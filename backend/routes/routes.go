package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/config"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/middleware"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/repositories"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

func ConfigRoutes(router *gin.Engine, db *pgxpool.Pool, cfg *config.Config, googleConfig, githubConfig *oauth2.Config, redisClient *redis.Client) {
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

	//CACHE CONTROL
	cache10Minutes := middleware.CacheControl(10 * time.Minute)
	//cache1Hour := middleware.CacheControl(1 * time.Hour)
	cache24Hours := middleware.CacheControl(24 * time.Hour)

	//User
	userRepo := repositories.NewUserRepository(db)

	// Algorithm Handlers and Services
	algoRepo := repositories.NewAlgorithmRepository(db)
	algoService := services.NewAlgorithmService(algoRepo, userRepo)
	algoHandler := handlers.NewAlgorithmHandler(algoService)

	//Auth
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, userRepo, redisClient, cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.AppName, cfg.EncryptSecretKey, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	authHandler := handlers.NewAuthHandler(authService, isProd, cfg.AppDomain, cfg.JwtRefreshExpiresDays)
	authSocialHandler := handlers.NewAuthSocialHandler(authService, googleConfig, githubConfig, cfg.AppDomain, cfg.FrontendUrl, isProd, cfg.JwtRefreshExpiresDays)

	//Sign up
	signUpService := services.NewSignUpService(userRepo, authRepo, *argonParams, cfg.JwtAccessSecret, cfg.JwtRefreshSecret, cfg.AppName, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	signUpHandler := handlers.NewSignUpHandler(signUpService, cfg.JwtRefreshExpiresDays, cfg.AppDomain, isProd)

	//TwoFactor
	twoFactorService := services.NewTwoFactorService(userRepo, cfg.EncryptSecretKey, cfg.AppName)
	twoFactorHandler := handlers.NewTwoFactorHandler(twoFactorService)

	//UserConfig
	emailService := services.NewEmailService(cfg.HostEmail, cfg.PortEmail, cfg.UserEmail, cfg.PasswordEmail, cfg.FromEmail, cfg.FrontendUrl, cfg.AppName)
	userConfigService := services.NewUserConfigService(userRepo, authRepo, *emailService, *argonParams, cfg.JwtRefreshSecret, cfg.AppName)
	userConfigHandler := handlers.NewUserConfigHandler(userConfigService)

	api := router.Group("/api")
	{
		sitemaps := api.Group("/sitemap", strictAbuseLimiter, cache24Hours, middleware.SitemapMiddleware(cfg.SitemapSecret))
		{
			sitemaps.GET("/algorithms", algoHandler.SitemapAlgorithms)
		}

		publicStandard := api.Group("", standardApiLimiter)
		{
			publicStandard.GET("/ping", handlers.AnswerPing)
			publicStandard.GET("/algorithms", cache10Minutes, algoHandler.ListAlgorithms)
			publicStandard.GET("/algorithms/:slugAndId", cache24Hours, algoHandler.GetAlgorithm)
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

			authenticatedAuth := auth.Group("", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName, redisClient), authFlowLimiter)
			{
				authenticatedAuth.POST("/logout", authHandler.Logout)
				authenticatedAuth.POST("/logout/all", authHandler.LogoutAll)
			}
		}

		users := api.Group("/users", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName, redisClient))
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

		admin := api.Group("/admin", middleware.AuthMiddleware(cfg.JwtAccessSecret, cfg.AppName, redisClient))
		{
			admin.Use(middleware.Fake404Middleware(cfg.AdminHash))
			admin.Use(middleware.EmployeeMiddleware())
			admin.Use(tenMbSize)
			admin.Use(standardApiLimiter)

			admin.GET("/ping", handlers.AnswerPing)

			algorithms := admin.Group("/algorithms")
			{
				create := algorithms.Group("", middleware.PermissionMiddleware("create:algorithms"))
				{
					create.GET("", algoHandler.GetAdminAlgorithms)
					create.POST("", algoHandler.PostAlgorithm)
					create.PUT("/:slugAndId", algoHandler.PutAlgorithm)
					create.GET("/:slugAndId", algoHandler.GetAdminAlgorithm)
					create.DELETE("/:slugAndId", algoHandler.DeleteAlgorithm)
					create.PATCH("/restore/:slugAndId", algoHandler.RestoreAlgorithm)
				}
			}

		}
	}
}
