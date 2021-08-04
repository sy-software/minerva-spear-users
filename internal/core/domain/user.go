package domain

import "time"

type User struct {
	Id string `json:"id,omitempty"`
	// User screen name, used for login
	Username string `json:"username,omitempty"`
	// User real name
	Name string `json:"name,omitempty"`
	// Optional url of the user display image
	Picture string `json:"picture,omitempty"`
}

type Login struct {
	// User screen name, used for login
	Username string `json:"username,omitempty"`
	// The OAuth2 provider used by this user
	Provider string `json:"provider,omitempty"`
	// The identifier connection this user with the OAuth provider
	TokenID string `json:"tokenID,omitempty"`
}

type Register struct {
	// User screen name, used for login
	Username string `json:"username,omitempty"`
	// User real name
	Name string `json:"name,omitempty"`
	// Optional url of the user display image
	Picture string `json:"picture,omitempty"`
	// For RBAC operations
	Role string `json:"role,omitempty"`
	// The OAuth2 provider used by this user
	Provider string `json:"provider,omitempty"`
	// The identifier connection this user with the OAuth provider
	TokenID string `json:"tokenID,omitempty"`
}

type UserToken struct {
	// The JWT for other requests authentication
	AccessToken string `json:"accessToken"`
	// To create a new JWT without re-enter credentials
	RefreshToken string `json:"refreshToken"`
	// How the token should be used I.E.: Bearer, Header, cookie, etc.
	TokenType string `json:"tokenType"`
	// When will this token expires
	ExpireTime time.Time `json:"expireTime"`
	// The user full info
	Info User `json:"info"`
}
