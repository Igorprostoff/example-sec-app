package usecase

import (

)

// UserUseCase -.
type UserUseCase struct {
	Repo   UserRepo
}

// New -.
func New(r UserRepo) *UserUseCase {
	return &UserUseCase{
		Repo:   r,
	}
}
