package storage

import (
	"context"
	"database/sql"

	"tgVideoCall/domain/admin/models"

	sq "github.com/Masterminds/squirrel"
)

// Пытается получить админа по userID
func (p *DB) GetAdmin(ctx context.Context, userID int) (models.Admin, error) {
	const op = "gates.postgres.GetAdmin"
	p.log.Debug(op, "Starting retrieving admin, for user: ", userID)

	query := p.sq.Select("user_id", "role").
		From("admins").
		Where(sq.Eq{"user_id": userID})

	qry, args, err := query.ToSql()
	p.log.Debug(op, "Building query: ", qry, "args: ", args)
	if err != nil {
		p.log.Error(op, "Error building query: ", err)
		return models.Admin{}, err
	}

	var result models.Admin
	err = p.db.QueryRowxContext(ctx, qry, args...).StructScan(&result)

	if err == sql.ErrNoRows {
		p.log.Debug(op, "Admin does not exist, for user: ", userID)
		return models.Admin{}, models.ErrNotAdmin
	}
	if err != nil {
		p.log.Error(op, "Error retrieving admin: ", err)
		return models.Admin{}, err
	}

	p.log.Debug(op, "Successfully retrieved admin for user: ", userID)
	return result, nil
}
