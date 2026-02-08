package admin

import (
	"context"

	"tgVideoCall/domain/admin/models"
)

type Storage interface {
	GetAdmin(ctx context.Context,userID int) (models.Admin, error)
}
