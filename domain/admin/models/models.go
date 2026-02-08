package models

import (
	"errors"
)

var (
	ErrNotAdmin = errors.New("user is not admin")
)

type Admin struct {
	UserID int64  `db:"user_id"`
	Role   string `db:"role"`
}
