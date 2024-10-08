package handlers

import (
	"encoding/json"
	"github.com/NastyaAR/music_library/internal/domain"
	"github.com/NastyaAR/music_library/internal/pkg/date_validate"
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

func getDateFromUser(date string) (time.Time, error) {
	if !date_validate.IsValidDate(date) {
		return time.Time{}, domain.ErrBadReleaseDate
	}

	parts := strings.Split(date, ".")

	numOfDate := make([]int, 3)
	var err error
	for i, part := range parts {
		numOfDate[i], err = strconv.Atoi(part)
		if err != nil {
			return time.Time{}, domain.ErrBadReleaseDate
		}
	}

	res := time.Date(numOfDate[2], time.Month(numOfDate[1]), numOfDate[0], 0, 0, 0, 0, time.UTC)
	return res, nil
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
		h.lg.Warn("song handler: create error: read", zap.Error(err))
		error_handler.NewError(ctx, domain.ErrInternalServer)
		return
	}
	err = json.Unmarshal(body, &songRequest)
	if err != nil {
		h.lg.Warn("song handler: create error: unmarsh", zap.Error(err), zap.Any("req", ctx.Request.Body))
		error_handler.NewError(ctx, domain.ErrInternalServer)
		return
	}

	var date time.Time

	if songRequest.ReleaseDate != "" {
		date, err = getDateFromUser(songRequest.ReleaseDate)
		if err != nil {
			h.lg.Warn("song handler: create error: unmarsh", zap.Error(err))
			error_handler.NewError(ctx, err)
			return
		}
	}

	song := domain.Song{
		Group:       songRequest.Group,
		Name:        songRequest.Name,
		ReleaseDate: date,
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

	var date time.Time

	if songRequest.ReleaseDate != "" {
		date, err = getDateFromUser(songRequest.ReleaseDate)
		if err != nil {
			h.lg.Warn("song handler: create error: unmarsh", zap.Error(err))
			error_handler.NewError(ctx, err)
			return
		}
	}

	song := domain.Song{
		Group:       songRequest.Group,
		Name:        songRequest.Name,
		ReleaseDate: date,
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

	var date time.Time

	if songRequest.ReleaseDate != "" {
		date, err = getDateFromUser(songRequest.ReleaseDate)
		if err != nil {
			h.lg.Warn("song handler: create error: unmarsh", zap.Error(err))
			error_handler.NewError(ctx, err)
			return
		}
	}

	filter := domain.Song{
		Group:       songRequest.Group,
		Name:        songRequest.Name,
		ReleaseDate: date,
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

// Get godoc
// @Summary      Get song info
// @Description  get song info
// @Tags         songs
// @Accept 		 json
// @Produce      json
// @Param        group    query     string  false  "group of song"
// @Param        name    query     string  false  "name of song"
// @Success      200  {object}  domain.GetSongResponse
// @Failure      400  {object}  error_handler.HTTPError
// @Failure      500  {object}  error_handler.HTTPError
// @Router       /info [get]
func (h *SongHandler) Get(ctx *gin.Context) {
	group := ctx.Request.URL.Query().Get("group")
	if group == "" {
		h.lg.Warn("song handler: get error")
		error_handler.NewError(ctx, domain.ErrBadGroup)
		return
	}

	name := ctx.Request.URL.Query().Get("name")
	if name == "" {
		h.lg.Warn("song handler: get error")
		error_handler.NewError(ctx, domain.ErrBadName)
		return
	}

	song, err := h.songUsecase.Get(ctx, group, name)
	if err != nil {
		h.lg.Warn("song handler: get error", zap.Error(err))
		error_handler.NewError(ctx, err)
		return
	}

	got := domain.GetSongResponse{
		ReleaseDate: getDate(song.ReleaseDate),
		Text:        song.Text,
		Link:        song.Link,
	}

	ctx.JSON(http.StatusOK, got)
}

// GetCouplet godoc
// @Summary      Get couplet with offset
// @Description  get couplet with offset
// @Tags         songs
// @Accept 		 json
// @Produce      json
// @Param        group    query     string  false  "group of song"
// @Param        name    query     string  false  "name of song"
// @Param        offset    query     string  false  "number of couplet"
// @Success      200  {object}  domain.GetCoupletResponse
// @Failure      400  {object}  error_handler.HTTPError
// @Failure      500  {object}  error_handler.HTTPError
// @Router       /songs/couplet [get]
func (h *SongHandler) GetCouplet(ctx *gin.Context) {
	group := ctx.Request.URL.Query().Get("group")
	if group == "" {
		h.lg.Warn("song handler: get error")
		error_handler.NewError(ctx, domain.ErrBadGroup)
		return
	}

	name := ctx.Request.URL.Query().Get("name")
	if name == "" {
		h.lg.Warn("song handler: get error")
		error_handler.NewError(ctx, domain.ErrBadName)
		return
	}

	offsetStr := ctx.Request.URL.Query().Get("offset")
	if offsetStr == "" {
		h.lg.Warn("song handler: getcouplet error")
		error_handler.NewError(ctx, domain.ErrBadOffset)
		return
	}

	offset, _ := strconv.Atoi(offsetStr)

	c, err := h.songUsecase.GetСouplet(ctx, group, name, offset)
	if err != nil {
		h.lg.Warn("song handler: get error", zap.Error(err))
		error_handler.NewError(ctx, err)
		return
	}

	couplet := domain.GetCoupletResponse{Couplet: c}

	ctx.JSON(http.StatusOK, couplet)
}
