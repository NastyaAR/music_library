package repo

import (
	"context"
	"fmt"
	"github.com/NastyaAR/music_library/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strings"
	"time"
)

type PostgresSongRepo struct {
	db *pgxpool.Pool
	lg *zap.Logger
}

func NewPostgresSongRepo(db *pgxpool.Pool, lg *zap.Logger) *PostgresSongRepo {
	lg.With(zap.String("component", "postgres_song_repo"))
	return &PostgresSongRepo{db: db, lg: lg}
}

func (p *PostgresSongRepo) Add(ctx context.Context, newSong *domain.Song) (domain.Song, error) {
	p.lg.Info("add new song", zap.Any("song", *newSong))

	query := `insert into songs(song_group, name, release_date, text, link)
	values ($1, $2, $3, $4, $5) returning *`

	var createdSong domain.Song
	err := p.db.QueryRow(ctx, query, newSong.Group, newSong.Name,
		newSong.ReleaseDate, newSong.Text, newSong.Link).Scan(&createdSong.Group, &createdSong.Name,
		&createdSong.ReleaseDate, &createdSong.Text, &createdSong.Link)
	if err != nil {
		p.lg.Warn("add error", zap.Error(err))
		return domain.Song{}, domain.ErrAddSongDB
	}

	p.lg.Info("successful adding new song")
	return createdSong, nil
}

func (p *PostgresSongRepo) Delete(ctx context.Context, group string, name string) error {
	p.lg.Info("delete song", zap.String("group", group),
		zap.String("name", name))

	query := `delete from songs where song_group=$1 and name=$2`
	_, err := p.db.Exec(ctx, query, group, name)
	if err != nil {
		p.lg.Warn("delete error", zap.Error(err))
		return domain.ErrDeleteSongDB
	}

	p.lg.Info("successful delete song")
	return nil
}

func (p *PostgresSongRepo) Update(ctx context.Context, group string, name string, upd *domain.Song) (domain.Song, error) {
	p.lg.Info("update song", zap.String("group", group),
		zap.String("name", name))

	query := `update songs set song_group=$1, name=$2, release_date=$3,
                 text=$4, link=$5
				where song_group=$6 and name=$7
				returning *`

	var newSong domain.Song
	err := p.db.QueryRow(ctx, query, upd.Group, upd.Name,
		upd.ReleaseDate, upd.Text, upd.Link, group, name).Scan(&newSong.Group, &newSong.Name,
		&newSong.ReleaseDate, &newSong.Text, &newSong.Link)
	if err != nil {
		p.lg.Warn("update error", zap.Error(err))
		return domain.Song{}, domain.ErrAddSongDB
	}

	p.lg.Info("successful updating new song")
	return newSong, nil
}

func (p *PostgresSongRepo) Get(ctx context.Context, group string, name string) (domain.Song, error) {
	p.lg.Info("get song", zap.String("group", group),
		zap.String("name", name))

	query := `select * from songs
	where song_group=$1 and name=$2`

	var newSong domain.Song
	err := p.db.QueryRow(ctx, query, group, name).Scan(&newSong.Group, &newSong.Name,
		&newSong.ReleaseDate, &newSong.Text, &newSong.Link)
	if err != nil {
		p.lg.Warn("get error", zap.Error(err))
		return domain.Song{}, domain.ErrAddSongDB
	}

	p.lg.Info("successful getting new song")
	return newSong, nil
}

func getFilterParams(filter *domain.Song) ([]string, []interface{}) {
	where := make([]string, 0)
	values := make([]interface{}, 0)
	cnt := 3
	if filter.Group != "" {
		where = append(where, fmt.Sprintf(`song_group=$%d`, cnt))
		cnt += 1
		values = append(values, filter.Group)
	}
	if filter.Name != "" {
		where = append(where, fmt.Sprintf(`name=$%d`, cnt))
		cnt += 1
		values = append(values, filter.Name)
	}

	var nilTime time.Time
	if filter.ReleaseDate != nilTime {
		where = append(where, fmt.Sprintf(`release_date=$%d`, cnt))
		cnt += 1
		values = append(values, filter.ReleaseDate)
	}

	if filter.Text != "" {
		where = append(where, fmt.Sprintf(`text=$%d`, cnt))
		cnt += 1
		values = append(values, filter.Text)
	}

	if filter.Link != "" {
		where = append(where, fmt.Sprintf(`link=$%d`, cnt))
		cnt += 1
		values = append(values, filter.Link)
	}

	return where, values
}

func (p *PostgresSongRepo) GetAll(ctx context.Context, filter *domain.Song, limit int, offset int) ([]domain.Song, error) {
	p.lg.Info("filter songs", zap.Any("filter", filter))

	query := `select * from songs`
	values := make([]interface{}, 0)
	values = append(values, limit, offset)
	where, critValues := getFilterParams(filter)

	if len(where) > 0 {
		query += ` where ` + strings.Join(where, ` and `)
	}

	values = append(values, critValues...)
	query += ` limit $1 offset $2`

	rows, err := p.db.Query(ctx, query, values...)
	defer rows.Close()
	if err != nil {
		p.lg.Warn("getall error", zap.Error(err))
		return nil, domain.ErrGetAllSongsDB
	}

	var song domain.Song
	songs := []domain.Song{}
	for rows.Next() {
		err = rows.Scan(&song.Group, &song.Name, &song.ReleaseDate,
			&song.Text, &song.Link)
		if err != nil {
			p.lg.Warn("getall error", zap.Error(err))
			continue
		}
		songs = append(songs, song)
	}

	return songs, nil
}
