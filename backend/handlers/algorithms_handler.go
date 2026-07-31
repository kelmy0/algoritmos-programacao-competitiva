package handlers

import (
	"log/slog"
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
	page, limit := parsePaginationQuery(c, 10)

	algorithms, finalPage, hasMore, err := h.Service.List(c.Request.Context(), page, limit)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ListAlgorithmsResponse{
		Page:       finalPage,
		Limit:      limit,
		HasMore:    hasMore,
		Algorithms: algorithms,
	})
}

func (h *AlgorithmHandler) GetAlgorithm(c *gin.Context) {
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

func (h *AlgorithmHandler) ListAdminAlgorithms(c *gin.Context) {
	userID, _, _, _, ok := GetAuthContext(c)
	if !ok {
		return
	}

	status := c.DefaultQuery("status", "")
	page, limit := parsePaginationQuery(c, 10)

	algorithms, finalPage, err := h.Service.ListAdmin(c.Request.Context(), page, limit, userID, status)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ListAlgorithmsResponse{
		Page:       finalPage,
		Limit:      limit,
		Algorithms: algorithms,
	})
}

func (h *AlgorithmHandler) GetAdminAlgorithm(c *gin.Context) {
	algoId, ok := parsePublicId(c)
	if !ok {
		return
	}

	userID, _, _, _, ok := GetAuthContext(c)
	if !ok {
		return
	}

	algorithm, err := h.Service.GetAdminAlgorithm(c.Request.Context(), algoId, userID)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.AlgorithmResponse{
		Data: algorithm,
	})
}

func (h *AlgorithmHandler) PostAlgorithm(c *gin.Context) {
	userID, _, _, _, ok := GetAuthContext(c)
	if !ok {
		return
	}

	var requestBody dto.PostAlgorithmRequest
	if !bindJSON(c, &requestBody) {
		return
	}

	algorithm, err := h.Service.PostAlgorithm(c.Request.Context(), requestBody, userID)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.AlgorithmResponse{
		Data: algorithm,
	})
}

func (h *AlgorithmHandler) DeleteAlgorithm(c *gin.Context) {
	userID, _, _, _, ok := GetAuthContext(c)
	if !ok {
		return
	}

	publicId, ok := parsePublicId(c)
	if !ok {
		return
	}

	err := h.Service.DeleteAlgorithm(c.Request.Context(), publicId, userID)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AlgorithmHandler) RestoreAlgorithm(c *gin.Context) {
	userID, _, _, _, ok := GetAuthContext(c)
	if !ok {
		return
	}

	publicId, ok := parsePublicId(c)
	if !ok {
		return
	}

	err := h.Service.RestoreAlgorithm(c.Request.Context(), publicId, userID)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AlgorithmHandler) PutAlgorithm(c *gin.Context) {
	userID, _, _, _, ok := GetAuthContext(c)
	if !ok {
		return
	}

	publicId, ok := parsePublicId(c)
	if !ok {
		return
	}

	var requestBody dto.PutAlgorithmRequest
	if !bindJSON(c, &requestBody) {
		return
	}

	algorithm, err := h.Service.PutAlgorithm(c.Request.Context(), requestBody, publicId, userID)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.AlgorithmResponse{
		Data: algorithm,
	})
}

func (h *AlgorithmHandler) SitemapAlgorithms(c *gin.Context) {
	algorithms, err := h.Service.SitemapAlgorithms(c.Request.Context())
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SitemapResponse{
		Data: algorithms,
	})
}

func (h *AlgorithmHandler) ListModerationAlgorithms(c *gin.Context) {
	userID, _, _, _, ok := GetAuthContext(c)
	if !ok {
		return
	}

	status := c.DefaultQuery("status", "")
	page, limit := parsePaginationQuery(c, 10)

	algorithms, finalPage, err := h.Service.ListModeration(c.Request.Context(), page, limit, userID, status)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ListAlgorithmsResponse{
		Page:       finalPage,
		Limit:      limit,
		Algorithms: algorithms,
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

func parsePaginationQuery(c *gin.Context, defaultLimit int) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	return page, limit
}

func bindJSON[T any](c *gin.Context, target *T) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		slog.Warn("bindJSON validation error", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return false
	}
	return true
}
