package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type AlgorithmHandler struct {
	Service *services.AlgorithmService
}

func NewAlgorithmHandler(service *services.AlgorithmService) *AlgorithmHandler {
	return &AlgorithmHandler{Service: service}
}

func (h *AlgorithmHandler) ListAlgorithms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	algorithms, finalPage, err := h.Service.List(c.Request.Context(), page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"ALGORITHMS_QUERY_FAILED",
			"Error querying algorithms.",
		))
		return
	}

	c.JSON(http.StatusOK, dto.ListAlgorithmsResponse{
		Page:  finalPage,
		Limit: limit,
		Data:  algorithms,
	})
}

func (h *AlgorithmHandler) GetAlgorithm(c *gin.Context) {
	//algorithm-slug-publicId
	publicId, ok := parsePublicId(c)
	if !ok {
		return
	}

	algorithm, err := h.Service.GetAlgorithmByPublicID(c.Request.Context(), publicId)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.AlgorithmResponse{
		Data: algorithm,
	})
}

func (h *AlgorithmHandler) PostAlgorithm(c *gin.Context) {
	var requestBody dto.PostAlgorithmRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	algorithm, err := h.Service.PostAlgorithm(c.Request.Context(), requestBody)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.AlgorithmResponse{
		Data: algorithm,
	})
}

func (h *AlgorithmHandler) DeleteAlgorithm(c *gin.Context) {
	publicId, ok := parsePublicId(c)
	if !ok {
		return
	}

	algorithm, err := h.Service.DeleteAlgorithm(c.Request.Context(), publicId)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.AlgorithmResponse{
		Data: algorithm,
	})
}

func (h *AlgorithmHandler) PutAlgorithm(c *gin.Context) {
	var requestBody dto.PutAlgorithmRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	algorithm, err := h.Service.PutAlgorithm(c.Request.Context(), requestBody)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.AlgorithmResponse{
		Data: algorithm,
	})
}

func parsePublicId(c *gin.Context) (string, bool) {
	slugAndId := c.Param("slugAndId")
	lastHifen := strings.LastIndex(slugAndId, "-")

	if slugAndId == "" || lastHifen == -1 {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_ALGORITHM_ID",
			"Invalid algorithm id format.",
		))
		return "", false
	}

	return slugAndId[lastHifen+1:], true
}
