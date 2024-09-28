package domain

import (
	"context"
	"errors"
	"time"
)

var ErrAddSongDB = errors.New("error while adding new song")
var ErrDeleteSongDB = errors.New("error while deleting song")
var ErrGetAllSongsDB = errors.New("error while getting songs")
var ErrNilCreateSongRequest = errors.New("bad nil request")
var ErrBadGroup = errors.New("bad group")
var ErrBadName = errors.New("bad name")
var ErrBadReleaseDate = errors.New("bad release date")
var ErrBadLimit = errors.New("bad limit")
var ErrBadOffset = errors.New("bad offset")
var ErrInternalServer = errors.New("something wrong while creating song")

var TimeLayout = "16.07.2006"

type Song struct {
	Group       string
	Name        string
	ReleaseDate time.Time
	Text        string
	Link        string
}

type UpdateSongRequest struct {
	Group       string `json:"group,omitempty"`
	Name        string `json:"name,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

type CreateSongResponse struct {
	Group       string `json:"group"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type CreateSongRequest struct {
	Group       string `json:"group"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

type GetSongResponse struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type GetSongsResponse struct {
	Songs []CreateSongResponse
}

type GetCoupletResponse struct {
	Couplet string `json:"couplet"`
}

type SongUsecase interface {
	Create(ctx context.Context, createReq *Song) (Song, error)
	Delete(ctx context.Context, group string, name string) error
	Update(ctx context.Context, group string, name string, updReq *Song) (Song, error)
	GetSongs(ctx context.Context, filter *Song, limit int, offset int) ([]Song, error)
	Get(ctx context.Context, group string, name string) (Song, error)
	GetСouplet(ctx context.Context, group string, name string, offset int) (string, error)
}

type SongRepo interface {
	Add(ctx context.Context, new *Song) (Song, error)
	Delete(ctx context.Context, group string, name string) error
	Update(ctx context.Context, group string, name string, upd *Song) (Song, error)
	Get(ctx context.Context, group string, name string) (Song, error)
	GetAll(ctx context.Context, filter *Song, limit int, offset int) ([]Song, error)
}
