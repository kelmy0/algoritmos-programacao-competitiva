package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

func HandleAPIError(c *gin.Context, err error) {
	if appErr, ok := errors.AsType[*models.AppError](err); ok {
		c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message))
		return
	}

	c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
		dto.CodeInternalError,
		dto.MsgUnexpectedError,
	))
}
