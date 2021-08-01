package mocks

import "github.com/sy-software/minerva-spear-users/internal/core/domain"

type UserRepo struct {
	CreateInterceptor        func(user domain.Register) (domain.User, error)
	GetByIdInterceptor       func(id string) (domain.User, error)
	GetByUsernameInterceptor func(username string) (domain.User, error)
}

func (repo *UserRepo) Create(user domain.Register) (domain.User, error) {
	return repo.CreateInterceptor(user)
}

func (repo *UserRepo) GetById(id string) (domain.User, error) {
	return repo.GetByIdInterceptor(id)
}

func (repo *UserRepo) GetByUsername(username string) (domain.User, error) {
	return repo.GetByUsernameInterceptor(username)
}
