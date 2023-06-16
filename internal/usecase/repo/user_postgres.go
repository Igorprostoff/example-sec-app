package repo

import (
	"context"
	"fmt"

	"securewebapp/internal/entity"
	"securewebapp/pkg/postgres"

	"github.com/jackc/pgx/v4"
	"github.com/microcosm-cc/bluemonday"
)

// UserRepo -.
type UserRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Login(ctx context.Context, user entity.User) (bool, error) {
	p := bluemonday.UGCPolicy()
	user.Email = p.Sanitize(user.Email)

	sql, args, err := r.Builder.Select("name").From("\"user\"").Where("password = $1 AND email = $2", user.Password, user.Email).ToSql()
	if err != nil {
		return false, fmt.Errorf("UserRepo - Login - r.Builder: %w", err)
	}

	var name string
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("UserRepo - Login - r.Pool.QueryRow: %w", err)
	}

	return true, nil
}

// Store -.
func (r *UserRepo) Store(ctx context.Context, u entity.User) error {
	sql, args, err := r.Builder.
		Insert("user").
		Columns("name, email, password").
		Values(u.Name, u.Email, u.Password).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}
