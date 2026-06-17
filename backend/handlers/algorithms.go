package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListAlgorithms(c *gin.Context) {
	algorithms := []string{"BFS", "DFS", "Dijkstra", "Segment tree"}

	c.JSON(http.StatusOK, gin.H{
		"data": algorithms,
	})
}
