package service

import (
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	mvdatetime "github.com/sy-software/minerva-go-utils/datetime"
	"github.com/sy-software/minerva-spear-users/internal/core/domain"
	"github.com/sy-software/minerva-spear-users/internal/core/ports"
)

const TOKEN_ISSUER = "minerva/spear/auth"
const TOKEN_AUDIENCE = "minerva/app"

type AuthService struct {
	repo   ports.UserRepo
	config domain.Config
}

func NewAuthService(repo ports.UserRepo, config domain.Config) *AuthService {
	return &AuthService{
		repo:   repo,
		config: config,
	}
}

// Creates a minerva JWT for a user validated by an OAuth provider
func (service *AuthService) Login(request domain.Login) (domain.UserToken, error) {
	return domain.UserToken{}, nil
}

// Registers a user validated by an OAuth provider into minerva platform
func (service *AuthService) Register(request domain.Register) (domain.UserToken, error) {
	newUser, err := service.repo.Create(request)

	if err != nil {
		return domain.UserToken{}, nil
	}
	now := mvdatetime.UnixUTCNow()
	expire := now.Add(time.Duration(service.config.Token.Duration) * time.Second)

	token := jwt.New()
	token.Set(jwt.IssuerKey, TOKEN_ISSUER)
	token.Set(jwt.ExpirationKey, expire)
	token.Set(jwt.SubjectKey, newUser.Id)
	token.Set(jwt.AudienceKey, TOKEN_AUDIENCE)

	token.Set("use", "access")

	refresh := jwt.New()
	refresh.Set(jwt.IssuerKey, TOKEN_ISSUER)
	refresh.Set(jwt.ExpirationKey, now.Add(time.Duration(service.config.Token.RefreshDuration)*time.Second))
	refresh.Set(jwt.SubjectKey, newUser.Id)
	refresh.Set(jwt.AudienceKey, TOKEN_AUDIENCE)

	refresh.Set("use", "refresh")

	key, err := service.config.Token.KeyPair()

	if err != nil {
		return domain.UserToken{}, err
	}

	serialized, err := jwt.Sign(token, jwa.RS256, key)

	if err != nil {
		return domain.UserToken{}, err
	}

	serializedRefresh, err := jwt.Sign(refresh, jwa.RS256, key)

	if err != nil {
		return domain.UserToken{}, err
	}

	return domain.UserToken{
		AccessToken:  string(serialized),
		RefreshToken: string(serializedRefresh),
		Info:         newUser,
		TokenType:    "Bearer",
		ExpireTime:   expire,
	}, nil
}

// Refresh the current user token
func (service *AuthService) Refresh(id string) (domain.UserToken, error) {
	return domain.UserToken{}, nil
}
