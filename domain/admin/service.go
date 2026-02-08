package admin

import (
	"context"
	"log/slog"

	"tgVideoCall/gates/storage"
	"tgVideoCall/domain/admin/models"
	"tgVideoCall/pkg/config"
)

type Service struct {
	log slog.Logger
	cfg config.Config
	ctx context.Context
	db  Storage
}

func NewService(ctx context.Context, cfg config.Config, log slog.Logger, db *storage.DB) *Service {
	return &Service{
		log: log,
		db:  db,
		cfg: cfg,
		ctx: ctx,
	}
}

func (s Service) GetAdmin(ctx context.Context, userID int) (models.Admin, error) {
	return s.db.GetAdmin(ctx, userID)
}

