package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/config"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
	authhandler "github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers/auth"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/middleware"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/repositories"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

func ConfigRoutes(router *gin.Engine, db *pgxpool.Pool, cfg *config.Config, googleConfig, githubConfig *oauth2.Config, redisClient *redis.Client) {
	isProd := cfg.AppEnv != "development"
	argonParams := utils.ArgonParams{
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
	issuerService := auth.NewSessionIssuer(authRepo, cfg.JwtAccessPrivateKey, cfg.JwtRefreshPrivateKey, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays, cfg.AppDomain)

	loginService := auth.NewLoginService(userRepo, issuerService, cfg.JwtAccessPrivateKey, cfg.AppDomain)
	loginHandler := authhandler.NewLoginHandler(loginService, cfg.JwtRefreshExpiresDays, isProd, cfg.AppDomain)

	authTwoFactorService := auth.NewTwoFactorService(userRepo, issuerService, redisClient, cfg.JwtAccessPublicKey, cfg.JwtAccessPrivateKey, cfg.AppDomain, cfg.EncryptSecretKey)
	authTwoFactorHandler := authhandler.NewTwoFactorHandler(authTwoFactorService, cfg.JwtRefreshExpiresDays, isProd, cfg.AppDomain)

	sessionService := auth.NewSessionService(authRepo, userRepo, redisClient, cfg.JwtRefreshPublicKey, cfg.JwtAccessPrivateKey, cfg.JwtRefreshPrivateKey, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays, cfg.AppDomain)
	sessionHandler := authhandler.NewSessionHandler(sessionService, cfg.JwtRefreshExpiresDays, isProd, cfg.AppDomain)

	googleProvider := authhandler.NewGoogleProvider(googleConfig)
	githubProvider := authhandler.NewGithubProvider(githubConfig)

	authSocialService := auth.NewSocialService(userRepo, issuerService, cfg.JwtAccessPrivateKey, cfg.AppDomain)
	authSocialHandler := authhandler.NewAuthSocialHandler(authSocialService, cfg.AppDomain, cfg.FrontendUrl, isProd, cfg.JwtRefreshExpiresDays, googleProvider, githubProvider)

	//Sign up
	signUpService := auth.NewSignUpService(userRepo, authRepo, argonParams, cfg.JwtAccessPrivateKey, cfg.JwtRefreshPrivateKey, cfg.AppDomain, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	signUpHandler := authhandler.NewSignUpHandler(signUpService, cfg.JwtRefreshExpiresDays, isProd, cfg.AppDomain)

	//TwoFactor
	twoFactorService := services.NewTwoFactorService(userRepo, authRepo, redisClient, cfg.EncryptSecretKey, cfg.AppName, cfg.AppDomain, cfg.JwtAccessPrivateKey, cfg.JwtRefreshPrivateKey, cfg.JwtAccessPublicKey, cfg.JwtRefreshPublicKey, cfg.JwtAccessExpiresMinutes, cfg.JwtRefreshExpiresDays)
	twoFactorHandler := handlers.NewTwoFactorHandler(twoFactorService, isProd, cfg.AppDomain, cfg.JwtRefreshExpiresDays)

	//UserConfig
	emailService := services.NewEmailService(cfg.HostEmail, cfg.PortEmail, cfg.UserEmail, cfg.PasswordEmail, cfg.FromEmail, cfg.FrontendUrl, cfg.AppName)
	userConfigService := services.NewUserConfigService(userRepo, authRepo, *emailService, argonParams, cfg.JwtRefreshPublicKey, cfg.AppDomain)
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
				authFlow.POST("/login", requireCaptcha, loginHandler.Login)
				authFlow.POST("/refresh", sessionHandler.Refresh)
				authFlow.POST("/sign-up", requireCaptcha, signUpHandler.SignUp)

				authFlow.GET("/:provider", authSocialHandler.SocialLogin)
				authFlow.GET("/:provider/callback", authSocialHandler.SocialCallback)
			}

			authStrict := auth.Group("", strictAbuseLimiter)
			{
				authStrict.POST("/forgot-password", requireCaptcha, userConfigHandler.ForgotPassword)
				authStrict.POST("/reset-password", requireCaptcha, userConfigHandler.ResetPassword)
				authStrict.POST("/verify-2fa", requireCaptcha, authTwoFactorHandler.Verify2FA)
			}

			authenticatedAuth := auth.Group("", requireAuth, authFlowLimiter)
			{
				authenticatedAuth.POST("/logout", sessionHandler.Logout)
				authenticatedAuth.POST("/logout/others", sessionHandler.LogoutOtherDevices)
				authenticatedAuth.POST("/logout/all", sessionHandler.LogoutAllDevices)
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

				/*
					linkSocial := me.Group("/link-social")
					{
						linkSocial.GET("/google", authSocialHandler.GoogleLinkAccount)
						linkSocial.GET("/github", authSocialHandler.GithubLinkAccount)
					}
				*/

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
