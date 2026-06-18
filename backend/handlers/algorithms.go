package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type AlgorithmHandler struct {
	service *services.AlgorithmService
}

func NewAlgorithmHandler(service *services.AlgorithmService) *AlgorithmHandler {
	return &AlgorithmHandler{service: service}
}

// List algorithms
func (h *AlgorithmHandler) ListAlgorithms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	algorithms, finalPage, err := h.service.List(c.Request.Context(), page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying algorithms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  finalPage,
		"limit": limit,
		"data":  algorithms,
	})
}

// Get algorithm by id
func (h *AlgorithmHandler) GetAlgorithm(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}

	algorithm, err := h.service.GetAlgorithmById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error geting content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": algorithm,
	})
}
