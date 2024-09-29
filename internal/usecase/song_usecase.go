package usecase

import (
	"context"
	"fmt"
	"github.com/NastyaAR/music_library/internal/domain"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strings"
	"time"
)

type SongUsecase struct {
	songRepo  domain.SongRepo
	validate  *validator.Validate
	lg        *zap.Logger
	dbTimeout time.Duration
}

func NewSongUsecase(songRepo domain.SongRepo, valid *validator.Validate, lg *zap.Logger) *SongUsecase {
	lg.With(zap.String("component", "song usecase"))
	return &SongUsecase{
		songRepo:  songRepo,
		validate:  valid,
		lg:        lg,
		dbTimeout: time.Hour,
	}
}

func (s *SongUsecase) Create(ctx context.Context,
	createReq *domain.Song) (domain.Song, error) {
	s.lg.Info("create song", zap.Any("request", *createReq))

	if createReq == nil {
		s.lg.Warn("create error: nil request",
			zap.Error(domain.ErrNilCreateSongRequest))
		return domain.Song{}, domain.ErrNilCreateSongRequest
	}

	err := s.validate.Var(createReq.Group, "required")
	if err != nil {
		s.lg.Warn("create error: bad group",
			zap.Error(domain.ErrBadGroup))
		return domain.Song{}, domain.ErrBadGroup
	}

	err = s.validate.Var(createReq.Name, "required")
	if err != nil {
		s.lg.Warn("create error: bad name",
			zap.Error(domain.ErrBadName))
		return domain.Song{}, domain.ErrBadName
	}

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	created, err := s.songRepo.Add(dbCtx, createReq)
	if err != nil {
		s.lg.Warn("create error", zap.Error(err))
		return domain.Song{},
			fmt.Errorf("create error: %v", err.Error())
	}

	s.lg.Info("successful create song")
	return created, nil
}

func (s *SongUsecase) Delete(ctx context.Context, group string, name string) error {
	s.lg.Info("delete song", zap.String("group", group),
		zap.String("name", name))

	if group == "" {
		s.lg.Warn("delete error: bad group",
			zap.Error(domain.ErrBadGroup))
		return domain.ErrBadGroup
	}

	if name == "" {
		s.lg.Warn("delete error: bad name",
			zap.Error(domain.ErrBadName))
		return domain.ErrBadName
	}

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	err := s.songRepo.Delete(dbCtx, group, name)
	if err != nil {
		s.lg.Warn("delete error", zap.Error(err))
		return fmt.Errorf("delete error: %v", err.Error())
	}

	s.lg.Info("successful delete")
	return nil
}

func (s *SongUsecase) Update(ctx context.Context, group string, name string, updReq *domain.Song) (domain.Song, error) {
	s.lg.Info("update song", zap.Any("request", *updReq))

	if group == "" {
		s.lg.Warn("update error: bad group",
			zap.Error(domain.ErrBadGroup))
		return domain.Song{}, domain.ErrBadGroup
	}

	if name == "" {
		s.lg.Warn("update error: bad name",
			zap.Error(domain.ErrBadName))
		return domain.Song{}, domain.ErrBadName
	}

	if updReq == nil {
		s.lg.Warn("update error: nil request",
			zap.Error(domain.ErrNilCreateSongRequest))
		return domain.Song{}, domain.ErrNilCreateSongRequest
	}

	err := s.validate.Var(updReq.Group, "required")
	if err != nil {
		s.lg.Warn("update error: bad group",
			zap.Error(domain.ErrBadGroup))
		return domain.Song{}, domain.ErrBadGroup
	}

	err = s.validate.Var(updReq.Name, "required")
	if err != nil {
		s.lg.Warn("update error: bad name",
			zap.Error(domain.ErrBadName))
		return domain.Song{}, domain.ErrBadName
	}

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	updated, err := s.songRepo.Update(dbCtx, group, name, updReq)
	if err != nil {
		s.lg.Warn("update error", zap.Error(err))
		return domain.Song{}, fmt.Errorf("update error: %v", err.Error())
	}

	s.lg.Info("successful update")
	return updated, nil
}

func (s *SongUsecase) GetSongs(ctx context.Context, filter *domain.Song,
	limit int, offset int) ([]domain.Song, error) {
	s.lg.Info("get songs", zap.Any("filter", *filter))

	if filter == nil {
		s.lg.Warn("getsongs error: nil request",
			zap.Error(domain.ErrNilCreateSongRequest))
		return nil, domain.ErrNilCreateSongRequest
	}

	if limit <= 0 {
		s.lg.Warn("getsongs error: bad limit",
			zap.Error(domain.ErrBadLimit))
		return nil, domain.ErrBadLimit
	}

	if offset < 1 {
		s.lg.Warn("getsongs error: bad offset",
			zap.Error(domain.ErrBadOffset))
		return nil, domain.ErrBadOffset
	}

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	songs, err := s.songRepo.GetAll(dbCtx, filter, limit, offset-1)
	if err != nil {
		s.lg.Warn("getsongs error", zap.Error(err))
		return nil, fmt.Errorf("getsongs error: %v", err.Error())
	}

	s.lg.Info("successful getsongs")
	return songs, nil
}

func (s *SongUsecase) Get(ctx context.Context, group string, name string) (domain.Song, error) {
	s.lg.Info("get", zap.String("group", group),
		zap.String("name", name))

	if group == "" {
		s.lg.Warn("get error: bad group",
			zap.Error(domain.ErrBadGroup))
		return domain.Song{}, domain.ErrBadGroup
	}

	if name == "" {
		s.lg.Warn("get error: bad name",
			zap.Error(domain.ErrBadName))
		return domain.Song{}, domain.ErrBadName
	}

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	song, err := s.songRepo.Get(dbCtx, group, name)
	if err != nil {
		s.lg.Warn("get error", zap.Error(err))
		return domain.Song{}, fmt.Errorf("get error: %v", err.Error())
	}

	s.lg.Info("successful get song")
	return song, nil
}

func (s *SongUsecase) GetСouplet(ctx context.Context,
	group string, name string, offset int) (string, error) {
	s.lg.Info("getcouplet", zap.String("group", group),
		zap.String("name", name))

	if group == "" {
		s.lg.Warn("getcouplet error: bad group",
			zap.Error(domain.ErrBadGroup))
		return "", domain.ErrBadGroup
	}

	if name == "" {
		s.lg.Warn("getcouplet error: bad name",
			zap.Error(domain.ErrBadName))
		return "", domain.ErrBadName
	}

	if offset < 1 {
		s.lg.Warn("getcouplet error: bad offset",
			zap.Error(domain.ErrBadOffset))
		return "", domain.ErrBadOffset
	}

	dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()

	song, err := s.songRepo.Get(dbCtx, group, name)
	if err != nil {
		s.lg.Warn("getcouplet error", zap.Error(err))
		return "", fmt.Errorf("getcouplet error: %v", err.Error())
	}

	couplets := strings.Split(song.Text, "\n\n")

	if offset-1 >= len(couplets) {
		s.lg.Warn("getcouplet error", zap.Error(err))
		return "", fmt.Errorf("getcouplet error: %v", domain.ErrBadOffset)
	}

	return couplets[offset-1], nil
}
