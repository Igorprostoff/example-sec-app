// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"securewebapp/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	UserRepo interface {
		Login(context.Context, entity.User) (bool, error)
		Store(context.Context, entity.User) error
	}
)
