package app

import (
	"fmt"
	"github.com/NastyaAR/music_library/internal/config"
	"github.com/NastyaAR/music_library/internal/delivery/http/v1/handlers"
	pkg "github.com/NastyaAR/music_library/internal/pkg/logger"
	repo "github.com/NastyaAR/music_library/internal/repo/postgres"
	"github.com/NastyaAR/music_library/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattes/migrate/source/file"
	"golang.org/x/net/context"
	"log"
	"time"
)

func Run() {
	cfg, err := config.ReadConfig("internal/config/config.yml")
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)

	logger, err := pkg.CreateLogger(cfg.LogFile, cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password,
		cfg.Host, cfg.Port, cfg.Db.Name)

	migr, err := migrate.New("file://migrations", connString)
	if err != nil {
		log.Fatalf("can't migrate postgresql: %v", err.Error())
	}
	err = migr.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("can't run up migrations on postgresql: %v", err.Error())
	}

	pool, err := pgxpool.New(ctx, connString)
	defer pool.Close()
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	songRepo := repo.NewPostgresSongRepo(pool, logger)
	validate := validator.New()
	songUsecase := usecase.NewSongUsecase(songRepo, validate, logger)

	songHandler := handlers.NewSongHandler(songUsecase, logger)
	router := gin.Default()
	router.POST("/songs", songHandler.Create)
	router.DELETE("/songs", songHandler.Delete)
	router.PATCH("/songs", songHandler.Update)
	router.GET("/songs", songHandler.GetSongs)
	router.GET("/info", songHandler.Get)
	router.GET("/songs/couplet", songHandler.GetCouplet)

	router.Run(":8080")
}
