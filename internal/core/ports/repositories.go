package ports

import "github.com/sy-software/minerva-spear-users/internal/core/domain"

type UserRepo interface {
	Create(user domain.Register) (domain.User, error)
	GetById(id string) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
}
