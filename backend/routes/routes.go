package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
)

func ConfigRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/ping", handlers.AnswerPing)
		api.GET("/list-algorithms", handlers.ListAlgorithms)
	}
}
