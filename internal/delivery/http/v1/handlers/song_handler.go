package handlers

import (
	"encoding/json"
	"github.com/NastyaAR/music_library/internal/domain"
	"github.com/NastyaAR/music_library/internal/pkg/error_handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type SongHandler struct {
	songUsecase domain.SongUsecase
	lg          *zap.Logger
}

func NewSongHandler(s domain.SongUsecase, lg *zap.Logger) *SongHandler {
	return &SongHandler{
		songUsecase: s,
		lg:          lg,
	}
}

// Create godoc
// @Summary      Create song
// @Description  create song
// @Tags         songs
// @Accept       json
// @Produce      json
// @Success      200  {object}  domain.Song
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /songs [post]
func (h *SongHandler) Create(ctx *gin.Context) {
	var songRequest domain.CreateSongRequest

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.lg.Warn("song handler: create error", zap.Error(err))
		error_handler.NewError(ctx, http.StatusInternalServerError,
			domain.ErrInternalServer)
		return
	}
	err = json.Unmarshal(body, &songRequest)
	if err != nil {
		h.lg.Warn("song handler: create error", zap.Error(err))
		error_handler.NewError(ctx, http.StatusInternalServerError,
			domain.ErrInternalServer)
		return
	}

	releaseDate, err := time.Parse(domain.TimeLayout, songRequest.ReleaseDate)
	if err != nil {
		h.lg.Warn("song handler: create error", zap.Error(err))
		error_handler.NewError(ctx, http.StatusInternalServerError,
			domain.ErrInternalServer)
		return
	}

	song := domain.Song{
		Group:       songRequest.Group,
		Name:        songRequest.Name,
		ReleaseDate: releaseDate,
		Text:        songRequest.Text,
		Link:        songRequest.Link,
	}

	created, err := h.songUsecase.Create(ctx, &song)
	if err != nil {
		h.lg.Warn("song handler: create error", zap.Error(err))
		error_handler.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, created)
}
