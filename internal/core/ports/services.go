package ports

import "github.com/sy-software/minerva-spear-users/internal/core/domain"

type UsersService interface {
	Login(request domain.Login) (domain.UserToken, error)
	Register(request domain.Login) (domain.UserToken, error)
	Refresh(id string) (domain.User, error)

	Get() (domain.User, error)
	GetById(id string) (domain.User, error)
}
