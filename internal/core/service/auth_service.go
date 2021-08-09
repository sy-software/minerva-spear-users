package service

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	mvdatetime "github.com/sy-software/minerva-go-utils/datetime"
	"github.com/sy-software/minerva-spear-users/internal/core/domain"
	"github.com/sy-software/minerva-spear-users/internal/core/ports"
)

const TOKEN_ISSUER = "minerva/spear/auth"
const TOKEN_AUDIENCE = "minerva/app"

type TokenUse string

const (
	Access  TokenUse = "access"
	Refresh TokenUse = "refresh"
)

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
// TODO: Implement token count limit
// TODO: Update user info on each new login
func (service *AuthService) Login(request domain.Login) (domain.UserToken, error) {
	user, err := service.repo.GetByUsername(request.Username)

	if err != nil {
		return domain.UserToken{}, err
	}

	key, err := service.config.Token.KeyPair()

	if err != nil {
		return domain.UserToken{}, err
	}

	return createUserToken(user, key, &service.config)
}

// Registers a user validated by an OAuth provider into minerva platform
func (service *AuthService) Register(request domain.Register) (domain.UserToken, error) {
	newUser, err := service.repo.Create(request)

	if err != nil {
		return domain.UserToken{}, err
	}

	now := mvdatetime.UnixUTCNow()
	expire := now.Add(time.Duration(service.config.Token.Duration) * time.Second)

	key, err := service.config.Token.KeyPair()

	if err != nil {
		return domain.UserToken{}, err
	}

	token, err := createToken(
		newUser.Id,
		expire,
		Access,
		&newUser,
		key,
	)

	if err != nil {
		return domain.UserToken{}, err
	}

	refresh, err := createToken(
		newUser.Id,
		now.Add(time.Duration(service.config.Token.RefreshDuration)*time.Second),
		Refresh,
		nil,
		key,
	)

	if err != nil {
		return domain.UserToken{}, err
	}

	if err != nil {
		return domain.UserToken{}, err
	}

	return domain.UserToken{
		AccessToken:  token,
		RefreshToken: refresh,
		Info:         newUser,
		TokenType:    "Bearer",
		ExpireTime:   expire,
	}, nil
}

// Refresh the current user token
// TODO: Implement single use refresh token, I.E.: Can't use same refresh token twice
func (service *AuthService) Refresh(refreshToken string) (domain.UserToken, error) {
	key, err := service.config.Token.KeyPair()

	if err != nil {
		return domain.UserToken{}, err
	}

	decoded, err := jwt.Parse(
		[]byte(refreshToken),
		jwt.WithVerify(jwa.RS256, key.PublicKey),
		jwt.WithValidate(true),
	)

	if err != nil {
		return domain.UserToken{}, err
	}

	use, ok := decoded.Get("use")
	if !ok || use.(string) != string(Refresh) {
		return domain.UserToken{}, errors.New("expected refresh token")
	}

	userId := decoded.Subject()

	user, err := service.repo.GetById(userId)

	if err != nil {
		return domain.UserToken{}, err
	}

	return createUserToken(user, key, &service.config)
}

// Get the current user information
func (service *AuthService) Me(userId string) (domain.User, error) {
	return service.repo.GetById(userId)
}

// Utils

func createUserToken(user domain.User, key *rsa.PrivateKey, config *domain.Config) (domain.UserToken, error) {
	now := mvdatetime.UnixUTCNow()
	expire := now.Add(time.Duration(config.Token.Duration) * time.Second)

	token, err := createToken(
		user.Id,
		expire,
		Access,
		&user,
		key,
	)

	if err != nil {
		return domain.UserToken{}, err
	}

	refresh, err := createToken(
		user.Id,
		now.Add(time.Duration(config.Token.RefreshDuration)*time.Second),
		Refresh,
		nil,
		key,
	)

	if err != nil {
		return domain.UserToken{}, err
	}

	if err != nil {
		return domain.UserToken{}, err
	}

	return domain.UserToken{
		AccessToken:  token,
		RefreshToken: refresh,
		Info:         user,
		TokenType:    "Bearer",
		ExpireTime:   expire,
	}, nil
}

func createToken(
	subject string,
	expire time.Time,
	use TokenUse,
	user *domain.User,
	key *rsa.PrivateKey,
) (string, error) {
	token := jwt.New()
	token.Set(jwt.IssuerKey, TOKEN_ISSUER)
	token.Set(jwt.ExpirationKey, expire)
	token.Set(jwt.SubjectKey, subject)
	token.Set(jwt.AudienceKey, TOKEN_AUDIENCE)

	token.Set("use", use)

	if user != nil {
		token.Set("user", user)
	}

	serialized, err := jwt.Sign(token, jwa.RS256, key)

	if err != nil {
		return "", nil
	}

	return string(serialized), nil
}
