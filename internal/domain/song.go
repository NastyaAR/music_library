package domain

import (
	"context"
	"time"
)

type Song struct {
	Group       string
	Name        string
	ReleaseDate time.Duration
	Text        string
	Link        string
}

type UpdateSongRequest struct {
	Group       string        `json:"group,omitempty"`
	Name        string        `json:"name,omitempty"`
	ReleaseDate time.Duration `json:"release_date,omitempty"`
	Text        string        `json:"text,omitempty"`
	Link        string        `json:"link,omitempty"`
}

type CreateSongResponse struct {
	Group       string        `json:"group"`
	Name        string        `json:"name"`
	ReleaseDate time.Duration `json:"release_date"`
	Text        string        `json:"text"`
	Link        string        `json:"link"`
}

type CreateSongRequest struct {
	Group       string        `json:"group"`
	Name        string        `json:"name"`
	ReleaseDate time.Duration `json:"release_date,omitempty"`
	Text        string        `json:"text,omitempty"`
	Link        string        `json:"link,omitempty"`
}

type GetSongResponse struct {
	ReleaseDate time.Duration `json:"release_date"`
	Text        string        `json:"text"`
	Link        string        `json:"link"`
}

type GetSongsResponse struct {
	Songs []CreateSongResponse
}

type SongUsecase interface {
	Create(ctx context.Context, createReq *CreateSongRequest) (CreateSongResponse, error)
	Delete(ctx context.Context, group string, name string) error
	Update(ctx context.Context, updReq *UpdateSongRequest) (CreateSongResponse, error)
	GetSongs(ctx context.Context, filter *UpdateSongRequest, limit int, offset int) (GetSongsResponse, error)
	Get(ctx context.Context, group string, name string) (GetSongResponse, error)
}

type SongRepo interface {
	Add(ctx context.Context, new *Song) (Song, error)
	Delete(ctx context.Context, group string, name string) error
	Update(ctx context.Context, upd *Song) (Song, error)
	Get(ctx context.Context, group string, name string) (Song, error)
	GetAll(ctx context.Context, filter *Song) ([]Song, error)
}
