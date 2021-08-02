package ports

import "github.com/sy-software/minerva-spear-users/internal/core/domain"

// AuthService handle all action related to user life cycle
type AuthService interface {
	// Creates a minerva JWT for a user validated by an OAuth provider
	Login(request domain.Login) (domain.UserToken, error)
	// Registers a user validated by an OAuth provider into minerva platform
	Register(request domain.Register) (domain.UserToken, error)
	// Refresh the current user token
	Refresh(refreshToken string) (domain.UserToken, error)
	// Get the current user information
	Me(userId string) (domain.User, error)
}
