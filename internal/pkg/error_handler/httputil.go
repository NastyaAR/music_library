package error_handler

import (
	"errors"
	"github.com/NastyaAR/music_library/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func isBadRequest(err error) bool {
	errorsList := []error{
		domain.ErrBadGroup,
		domain.ErrBadOffset,
		domain.ErrBadLimit,
		domain.ErrBadName,
		domain.ErrBadReleaseDate,
		domain.ErrQueryParams,
	}

	for _, e := range errorsList {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}

func NewError(ctx *gin.Context, err error) {
	var status int
	if isBadRequest(err) {
		status = http.StatusBadRequest
	} else {
		status = http.StatusInternalServerError
	}

	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}

	ctx.JSON(status, er)
}
