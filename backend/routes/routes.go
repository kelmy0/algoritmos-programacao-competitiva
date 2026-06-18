package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

func ConfigRoutes(router *gin.Engine, db *pgxpool.Pool) {
	algoService := services.NewAlgorithmService(db)
	algoHandler := handlers.NewAlgorithmHandler(algoService)

	api := router.Group("/api")
	{
		api.GET("/ping", handlers.AnswerPing)
		api.GET("/algorithms", algoHandler.ListAlgorithms)
		api.GET("/algorithms/:id", algoHandler.GetAlgorithm)
	}
}
