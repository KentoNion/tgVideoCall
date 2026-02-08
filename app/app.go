package app

import (
	"context"
	"fmt"
	"log/slog"

	"tgVideoCall/domain/admin"
	"tgVideoCall/gates/storage"
	"tgVideoCall/gates/telegram"
	"tgVideoCall/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	goose "github.com/pressly/goose/v3"
)

func Run(ctx context.Context, log slog.Logger, cfg config.Config) {
	const op = "app.Run"
	log.Debug("Debug logging on")

	//подключение к дб
	connstr := fmt.Sprintf("user=%s password=%s dbname=youtube_hub_bot host=%s sslmode=%s client_encoding=UTF8", cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Ssl)
	conn, err := sqlx.Connect("postgres", connstr) //драйвер и имя бд
	if err != nil {
		panic(err)
	}

	db := storage.NewPostgresDB(conn, log)

	//накатываем миграцию
	if err = goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	//goose.Down(conn.DB, cfg.DB.MigrationsPath)
	if err = goose.Up(conn.DB, cfg.DB.MigrationsPath); err != nil {
		panic(err)
	}

	service := admin.NewService(ctx, cfg, log, db)

	server := telegram.NewServer(ctx, log, cfg, *service)
	server.RunServer()
}
