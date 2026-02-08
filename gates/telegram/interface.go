package telegram

import (
	"context"

	"tgVideoCall/domain/admin/models"
)

type Service interface {
	// Получает администратора из базы данных, возвращает ошибку если не является администратором
	GetAdmin(ctx context.Context, id int)(admin models.Admin, err error)
}
