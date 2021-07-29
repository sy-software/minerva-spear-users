package domain

import "time"

type User struct {
	Id string `bson:"_id,omitempty" json:"id,omitempty"`
	// User screen name, used for login
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	// User real name
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	// Optional url of the user display image
	Picture string `bson:"picture,omitempty" json:"picture,omitempty"`
}

type Login struct {
	// User screen name, used for login
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	// For RBAC operations
	Role string `bson:"role,omitempty" json:"role,omitempty"`
	// The OAuth2 provider used by this user
	Provider string `bson:"provider,omitempty" json:"provider,omitempty"`
	// The identifier connection this user with the OAuth provider
	TokenID string `bson:"tokenID,omitempty" json:"tokenID,omitempty"`
}

type Register struct {
	// User screen name, used for login
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	// User real name
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	// Optional url of the user display image
	Picture string `bson:"picture,omitempty" json:"picture,omitempty"`
	// For RBAC operations
	Role string `bson:"role,omitempty" json:"role,omitempty"`
	// The OAuth2 provider used by this user
	Provider string `bson:"provider,omitempty" json:"provider,omitempty"`
	// The identifier connection this user with the OAuth provider
	TokenID string `bson:"tokenID,omitempty" json:"tokenID,omitempty"`
}

type UserToken struct {
	// The JWT for other requests authentication
	AccessToken string `bson:"accessToken" json:"accessToken"`
	// To create a new JWT without re-enter credentials
	RefreshToken string `bson:"refreshToken" json:"refreshToken"`
	// How the token should be used I.E.: Bearer, Header, cookie, etc.
	TokenType string `bson:"tokenType" json:"tokenType"`
	// When will this token expires
	ExpireTime time.Time `bson:"expireTime" json:"expireTime"`
	// The user full info
	Info User `bson:"info" json:"info"`
}
