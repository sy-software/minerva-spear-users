package ports

import "github.com/sy-software/minerva-spear-users/internal/core/domain"

// UserRepo handles interaction with user related storage operations
type UserRepo interface {
	// Create saves a new user
	Create(user domain.Register) (domain.User, error)
	// GetById looks for a user with the provided ID
	GetById(id string) (domain.User, error)
	// GetByUsernamelooks for a user with the provided username
	GetByUsername(username string) (domain.User, error)
}

// ConfigRepository provides connection to our config server
type ConfigRepository interface {
	// Get connects to the configuration server and loads the config
	Get() (domain.Config, error)
}
