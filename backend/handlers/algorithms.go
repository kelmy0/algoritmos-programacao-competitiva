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

func (h *AlgorithmHandler) ListAlgorithms(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	algorithms, err := h.service.List(limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying algorithms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"data":  algorithms,
	})
}

func (h *AlgorithmHandler) GetAlgorithm(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
	}

	algorithm, err := h.service.GetById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error geting content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": algorithm,
	})
}
