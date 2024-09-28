package handlers

import (
	"encoding/json"
	"github.com/NastyaAR/music_library/internal/domain"
	"github.com/NastyaAR/music_library/internal/pkg/error_handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func getDate(timestamp time.Time) string {
	date := strings.Split(timestamp.String(), " ")
	parts := strings.Split(date[0], "-")

	result := strings.Join([]string{parts[2], parts[1], parts[0]}, ".")
	return result
}

// Create godoc
// @Summary      Create song
// @Description  create song
// @Tags         songs
// @Accept       json
// @Produce      json
// @Success      200  {object}  domain.Song
// @Failure      400  {object}  error_handler.HTTPError
// @Failure      500  {object}  error_handler.HTTPError
// @Router       /songs [post]
func (h *SongHandler) Create(ctx *gin.Context) {
	var songRequest domain.CreateSongRequest

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.lg.Warn("song handler: create error", zap.Error(err))
		error_handler.NewError(ctx, domain.ErrInternalServer)
		return
	}
	err = json.Unmarshal(body, &songRequest)
	if err != nil {
		h.lg.Warn("song handler: create error", zap.Error(err))
		error_handler.NewError(ctx, domain.ErrInternalServer)
		return
	}

	var releaseDate time.Time

	if songRequest.ReleaseDate != "" {
		releaseDate, err = time.Parse(domain.TimeLayout, songRequest.ReleaseDate)
		if err != nil {
			h.lg.Warn("song handler: create error", zap.Error(err))
			error_handler.NewError(ctx, domain.ErrInternalServer)
			return
		}
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
		error_handler.NewError(ctx, err)
		return
	}

	createdResponse := domain.CreateSongResponse{
		Group:       created.Group,
		Name:        created.Name,
		ReleaseDate: getDate(created.ReleaseDate),
		Text:        created.Text,
		Link:        created.Link,
	}

	ctx.JSON(http.StatusOK, createdResponse)
}

// Delete godoc
// @Summary      Delete song
// @Description  delete song
// @Tags         songs
// @Produce      json
// @Param        group    query     string  false  "group of song"
// @Param        name    query     string  false  "name of song"
// @Success      200
// @Failure      400  {object}  error_handler.HTTPError
// @Failure      500  {object}  error_handler.HTTPError
// @Router       /songs [delete]
func (h *SongHandler) Delete(ctx *gin.Context) {
	group := ctx.Request.URL.Query().Get("group")
	if group == "" {
		h.lg.Warn("song handler: delete error")
		error_handler.NewError(ctx, domain.ErrBadGroup)
		return
	}

	name := ctx.Request.URL.Query().Get("name")
	if name == "" {
		h.lg.Warn("song handler: delete error")
		error_handler.NewError(ctx, domain.ErrBadName)
		return
	}

	err := h.songUsecase.Delete(ctx, group, name)
	if err != nil {
		h.lg.Warn("song handler: delete error", zap.Error(err))
		error_handler.NewError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// Update godoc
// @Summary      Update song
// @Description  update song
// @Tags         songs
// @Accept 		 json
// @Produce      json
// @Param        group    query     string  false  "group of song"
// @Param        name    query     string  false  "name of song"
// @Success      200  {object}  domain.Song
// @Failure      400  {object}  error_handler.HTTPError
// @Failure      500  {object}  error_handler.HTTPError
// @Router       /songs [patch]
func (h *SongHandler) Update(ctx *gin.Context) {
	group := ctx.Request.URL.Query().Get("group")
	if group == "" {
		h.lg.Warn("song handler: update error")
		error_handler.NewError(ctx, domain.ErrBadGroup)
		return
	}

	name := ctx.Request.URL.Query().Get("name")
	if name == "" {
		h.lg.Warn("song handler: update error")
		error_handler.NewError(ctx, domain.ErrBadName)
		return
	}

	var songRequest domain.UpdateSongRequest

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.lg.Warn("song handler: update error", zap.Error(err))
		error_handler.NewError(ctx, domain.ErrInternalServer)
		return
	}
	err = json.Unmarshal(body, &songRequest)
	if err != nil {
		h.lg.Warn("song handler: update error", zap.Error(err))
		error_handler.NewError(ctx, domain.ErrInternalServer)
		return
	}

	var releaseDate time.Time

	if songRequest.ReleaseDate != "" {
		releaseDate, err = time.Parse(domain.TimeLayout, songRequest.ReleaseDate)
		if err != nil {
			h.lg.Warn("song handler: create error", zap.Error(err))
			error_handler.NewError(ctx, domain.ErrInternalServer)
			return
		}
	}

	song := domain.Song{
		Group:       songRequest.Group,
		Name:        songRequest.Name,
		ReleaseDate: releaseDate,
		Text:        songRequest.Text,
		Link:        songRequest.Link,
	}

	updated, err := h.songUsecase.Update(ctx, group, name, &song)
	if err != nil {
		h.lg.Warn("song handler: update error", zap.Error(err))
		error_handler.NewError(ctx, err)
		return
	}

	createdResponse := domain.CreateSongResponse{
		Group:       updated.Group,
		Name:        updated.Name,
		ReleaseDate: getDate(updated.ReleaseDate),
		Text:        updated.Text,
		Link:        updated.Link,
	}

	ctx.JSON(http.StatusOK, createdResponse)
}

// GetSongs godoc
// @Summary      Get songs with filter, limit and offset
// @Description  get songs with filter, limit and offset
// @Tags         songs
// @Accept 		 json
// @Produce      json
// @Param        limit    query     string  false  "songs on page"
// @Param        offset    query     string  false  "page"
// @Success      200  {array}  domain.Song
// @Failure      400  {object}  error_handler.HTTPError
// @Failure      500  {object}  error_handler.HTTPError
// @Router       /songs [get]
func (h *SongHandler) GetSongs(ctx *gin.Context) {
	limitStr := ctx.Request.URL.Query().Get("limit")
	if limitStr == "" {
		h.lg.Warn("song handler: getsongs error")
		error_handler.NewError(ctx, domain.ErrBadLimit)
		return
	}

	limit, _ := strconv.Atoi(limitStr)

	offsetStr := ctx.Request.URL.Query().Get("offset")
	if offsetStr == "" {
		h.lg.Warn("song handler: getsongs error")
		error_handler.NewError(ctx, domain.ErrBadOffset)
		return
	}

	offset, _ := strconv.Atoi(offsetStr)

	var songRequest domain.UpdateSongRequest

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.lg.Warn("song handler: getsongs error", zap.Error(err))
		error_handler.NewError(ctx, domain.ErrInternalServer)
		return
	}
	err = json.Unmarshal(body, &songRequest)
	if err != nil {
		h.lg.Warn("song handler: getsongs error", zap.Error(err))
		error_handler.NewError(ctx, domain.ErrInternalServer)
		return
	}

	var releaseDate time.Time

	if songRequest.ReleaseDate != "" {
		releaseDate, err = time.Parse(domain.TimeLayout, songRequest.ReleaseDate)
		if err != nil {
			h.lg.Warn("song handler: getsongs error", zap.Error(err))
			error_handler.NewError(ctx, domain.ErrInternalServer)
			return
		}
	}

	filter := domain.Song{
		Group:       songRequest.Group,
		Name:        songRequest.Name,
		ReleaseDate: releaseDate,
		Text:        songRequest.Text,
		Link:        songRequest.Link,
	}

	songs, err := h.songUsecase.GetSongs(ctx, &filter,
		limit, offset)
	if err != nil {
		h.lg.Warn("song handler: getsongs error", zap.Error(err))
		error_handler.NewError(ctx, err)
		return
	}

	songsResponse := domain.GetSongsResponse{Songs: make([]domain.CreateSongResponse, 0)}
	for _, s := range songs {
		song := domain.CreateSongResponse{
			Group:       s.Group,
			Name:        s.Name,
			ReleaseDate: getDate(s.ReleaseDate),
			Text:        s.Text,
			Link:        s.Link,
		}
		songsResponse.Songs = append(songsResponse.Songs, song)
	}

	ctx.JSON(http.StatusOK, songsResponse)
}
