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
	isProd := cfg.AppEnv != "development"
	argonParams := &utils.ArgonParams{
		Memory:      cfg.Memory,
		Iterations:  cfg.Iterations,
		Parallelism: cfg.Parallelism,
		SaltLength:  cfg.SaltLength,
		KeyLength:   cfg.KeyLength,
	}
	//BODY SIZE
	hundredKbSize := middleware.LimitBodySize(1024 * 128)
	tenMbSize := middleware.LimitBodySize(10 * 1024 * 1024)

	//Query Limit
	fiveHundredQuerySize := middleware.LimitQueryParamsSize(512)
	thousandQuerySize := middleware.LimitQueryParamsSize(1024)
	twoThousandUrlSize := middleware.LimitQueryParamsSize(2048)

	//RATE LIMIT
	standardApiLimiter := middleware.RateLimitMiddleware(middleware.NewRateLimiter(redisClient, rate.Limit(5), 5))
	authFlowLimiter := middleware.RateLimitMiddleware(middleware.NewRateLimiter(redisClient, rate.Limit(0.1), 5))
	strictAbuseLimiter := middleware.RateLimitMiddleware(middleware.NewRateLimiter(redisClient, rate.Limit(0.0055), 2))

	//CACHE CONTROL
	cache10Minutes := middleware.CacheControl(10 * time.Minute)
	//cache1Hour := middleware.CacheControl(1 * time.Hour)
	cache24Hours := middleware.CacheControl(24 * time.Hour)

	//CAPTCHA
	requireCaptcha := middleware.RequireCaptcha(cfg.TurnstileSecret)

	//AUTH MIDDLEWARE
	requireAuth := middleware.AuthMiddleware(cfg.JwtAccessPublicKey, cfg.AppDomain, redisClient)

	//ADMIN MIDDLEWARES
	fake404 := middleware.Fake404Middleware(cfg.AdminHash)
	requireEmployee := middleware.EmployeeMiddleware()

	//User
	userRepo := repositories.NewUserRepository(db)

	// Algorithm Handlers and Services
	algoRepo := repositories.NewAlgorithmRepository(db)
	algoService := services.NewAlgorithmService(algoRepo, userRepo)
	algoHandler := handlers.NewAlgorithmHandler(algoService)

	//Auth
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, userRepo, redisClient, cfg.JwtAccessPrivateKey, cfg.JwtRefreshPrivateKey, cfg.JwtAccessPublicKey, cfg.JwtRefreshPublicKey, cfg.AppDomain, cfg.EncryptSecretKey, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	authHandler := handlers.NewAuthHandler(authService, isProd, cfg.AppDomain, cfg.JwtRefreshExpiresDays)
	authSocialHandler := handlers.NewAuthSocialHandler(authService, googleConfig, githubConfig, cfg.AppDomain, cfg.FrontendUrl, isProd, cfg.JwtRefreshExpiresDays)

	//Sign up
	signUpService := services.NewSignUpService(userRepo, authRepo, *argonParams, cfg.JwtAccessPrivateKey, cfg.JwtRefreshPrivateKey, cfg.AppDomain, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	signUpHandler := handlers.NewSignUpHandler(signUpService, cfg.JwtRefreshExpiresDays, cfg.AppDomain, isProd)

	//TwoFactor
	twoFactorService := services.NewTwoFactorService(userRepo, authRepo, redisClient, cfg.EncryptSecretKey, cfg.AppName, cfg.AppDomain, cfg.JwtAccessPrivateKey, cfg.JwtRefreshPrivateKey, cfg.JwtAccessPublicKey, cfg.JwtRefreshPublicKey, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	twoFactorHandler := handlers.NewTwoFactorHandler(twoFactorService, isProd, cfg.AppDomain, cfg.JwtRefreshExpiresDays)

	//UserConfig
	emailService := services.NewEmailService(cfg.HostEmail, cfg.PortEmail, cfg.UserEmail, cfg.PasswordEmail, cfg.FromEmail, cfg.FrontendUrl, cfg.AppName)
	userConfigService := services.NewUserConfigService(userRepo, authRepo, *emailService, *argonParams, cfg.JwtRefreshPublicKey, cfg.AppDomain)
	userConfigHandler := handlers.NewUserConfigHandler(userConfigService)

	api := router.Group("/api")
	{
		sitemaps := api.Group("/sitemap", fiveHundredQuerySize, strictAbuseLimiter, middleware.SitemapMiddleware(cfg.SitemapSecret), cache24Hours)
		{
			sitemaps.GET("/algorithms", algoHandler.SitemapAlgorithms)
		}

		publicStandard := api.Group("", standardApiLimiter)
		{
			publicStandard.GET("/ping", fiveHundredQuerySize, handlers.AnswerPing)
			publicStandard.GET("/algorithms", thousandQuerySize, cache10Minutes, algoHandler.ListAlgorithms)
			publicStandard.GET("/algorithms/:slugAndId", thousandQuerySize, cache24Hours, algoHandler.GetAlgorithm)
		}

		auth := api.Group("/auth", twoThousandUrlSize, hundredKbSize)
		{
			authFlow := auth.Group("", authFlowLimiter)
			{
				authFlow.POST("/login", requireCaptcha, authHandler.Auth)
				authFlow.POST("/refresh", authHandler.Refresh)
				authFlow.GET("/google", authSocialHandler.GoogleLogin)
				authFlow.GET("/google/callback", authSocialHandler.GoogleCallback)
				authFlow.GET("/github", authSocialHandler.GithubLogin)
				authFlow.GET("/github/callback", authSocialHandler.GithubCallback)
				authFlow.POST("/sign-up", requireCaptcha, signUpHandler.SignUp)
			}

			authStrict := auth.Group("", strictAbuseLimiter)
			{
				authStrict.POST("/forgot-password", requireCaptcha, userConfigHandler.ForgotPassword)
				authStrict.POST("/reset-password", requireCaptcha, userConfigHandler.ResetPassword)
				authStrict.POST("/verify-2fa", requireCaptcha, authHandler.Verify2FA)
			}

			authenticatedAuth := auth.Group("", requireAuth, authFlowLimiter)
			{
				authenticatedAuth.POST("/logout", authHandler.Logout)
				authenticatedAuth.POST("/logout/others", authHandler.LogoutOtherDevices)
				authenticatedAuth.POST("/logout/all", authHandler.LogoutAllDevices)
			}
		}

		users := api.Group("/users", twoThousandUrlSize, hundredKbSize)
		{
			me := users.Group("/me", requireAuth, authFlowLimiter)
			{
				me.GET("", userConfigHandler.GetMyCredentials)

				password := me.Group("/password")
				{
					password.POST("/set", userConfigHandler.DefinePassword)
					password.POST("/change", userConfigHandler.ChangePassword)
				}

				twoFa := me.Group("/2fa")
				{
					twoFa.POST("/generate", twoFactorHandler.Generate2FA)
					twoFa.POST("/enable", twoFactorHandler.Enable2FA)
					twoFa.POST("/disable", twoFactorHandler.Disable2FA)
				}

				linkSocial := me.Group("/link-social")
				{
					linkSocial.GET("/google", authSocialHandler.GoogleLinkAccount)
					linkSocial.GET("/github", authSocialHandler.GithubLinkAccount)
				}

			}
		}

		admin := api.Group("/admin", twoThousandUrlSize, tenMbSize, fake404, requireAuth, requireEmployee)
		{
			admin.GET("/ping", handlers.AnswerPing)

			algorithms := admin.Group("/algorithms")
			{
				createPerm := middleware.PermissionMiddleware("create:algorithms")

				standard := algorithms.Group("", standardApiLimiter)
				{
					standard.GET("", createPerm, algoHandler.ListAdminAlgorithms)
					standard.GET("/:slugAndId", createPerm, algoHandler.GetAdminAlgorithm)

					modPerm := middleware.PermissionMiddleware("moderate:algorithms")
					standard.GET("/moderation", modPerm, algoHandler.ListModerationAlgorithms)
				}

				authFlow := algorithms.Group("", authFlowLimiter, createPerm)
				{
					authFlow.POST("", algoHandler.PostAlgorithm)
					authFlow.PUT("/:slugAndId", algoHandler.PutAlgorithm)
					authFlow.DELETE("/:slugAndId", algoHandler.DeleteAlgorithm)
					authFlow.PATCH("/restore/:slugAndId", algoHandler.RestoreAlgorithm)
				}
			}
		}
	}
}
