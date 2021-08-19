package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sy-software/minerva-spear-users/internal/core/domain"
	"github.com/sy-software/minerva-spear-users/internal/core/ports"
)

const BEARER_REGEX = "^Bearer (.+)$"

const (
	REQUEST_ID_HEADER string = "X-REQUEST-ID"
	USER_INFO_HEADER  string = "X-USER-INFO"
	TOKEN_USE_HEADER  string = "X-TOKEN-USE"
	USER_ID_HEADER    string = "X-USER-ID"
)

type AuthRESTHandler struct {
	config  *domain.Config
	service ports.AuthService
}

func NewAuthRESTHandler(config *domain.Config, service ports.AuthService) *AuthRESTHandler {
	return &AuthRESTHandler{
		config:  config,
		service: service,
	}
}

func (handler *AuthRESTHandler) CreateRoutes(router *gin.Engine) {
	group := router.Group(handler.config.APIPrefix)
	{
		group.POST("/login", func(c *gin.Context) {
			token, err := handler.Login(c)

			if err != nil {
				handleError(err, c)
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": token})
		})

		group.POST("/refresh", func(c *gin.Context) {
			token, err := handler.Refresh(c)

			if err != nil {
				handleError(err, c)
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": token})
		})

		group.POST("/register", func(c *gin.Context) {
			token, err := handler.Register(c)

			if err != nil {
				handleError(err, c)
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": token})
		})

		group.POST("/authenticate", func(c *gin.Context) {
			token, err := handler.Authenticate(c)

			if err != nil {
				handleError(err, c)
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": token})
		})

		group.GET("/me", func(c *gin.Context) {
			user, err := handler.Me(c)

			if err != nil {
				handleError(err, c)
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": user})
		})
	}
}

func (handler *AuthRESTHandler) Login(c *gin.Context) (domain.UserToken, error) {
	var login domain.Login

	log.Debug().Msgf("Request Headers: %+v", c.Request.Header)
	userInfo := c.Request.Header.Get(USER_INFO_HEADER)

	if len(userInfo) == 0 {
		log.Error().Msg("User info header is not present")
		return domain.UserToken{}, &InvalidRequestError
	}

	userDecoded, err := base64.StdEncoding.DecodeString(userInfo)

	if err != nil {
		log.Error().Err(err).Msg("User info header can't be decoded from base64")
		return domain.UserToken{}, &InvalidRequestError
	}

	err = json.Unmarshal(userDecoded, &login)

	if err != nil {
		log.Error().Err(err).Msg("User info header can't be decoded from JSON")
		return domain.UserToken{}, &InvalidRequestError
	}

	user, err := handler.service.Login(login)

	if err != nil {
		log.Error().Err(err).Msg("Login error")
		if err.Error() == "not_found" {
			return domain.UserToken{}, &UserNotRegisteredErr
		}

		// Any other error is considered an unknown or unexpected error
		// user should only get internal server error
		return domain.UserToken{}, &InternalServerError
	}

	return user, nil
}

func (handler *AuthRESTHandler) Register(c *gin.Context) (domain.UserToken, error) {
	var register domain.Register

	userInfo := c.Request.Header.Get(USER_INFO_HEADER)

	if len(userInfo) == 0 {
		return domain.UserToken{}, &InvalidRequestError
	}

	userDecoded, err := base64.StdEncoding.DecodeString(userInfo)

	if err != nil {
		return domain.UserToken{}, &InvalidRequestError
	}

	err = json.Unmarshal(userDecoded, &register)

	if err != nil {
		return domain.UserToken{}, &InvalidRequestError
	}

	user, err := handler.service.Register(register)

	if err != nil {
		if err.Error() == "duplicated_value" {
			return domain.UserToken{}, &UserAlreadyRegisteredErr
		}

		// Any other error is considered an unknown or unexpected error
		// user should only get internal server error
		return domain.UserToken{}, &InternalServerError
	}

	return user, nil
}

func (handler *AuthRESTHandler) Refresh(c *gin.Context) (domain.UserToken, error) {
	refreshToken := c.GetHeader("Authorization")

	re := regexp.MustCompile(BEARER_REGEX)
	if !re.MatchString(refreshToken) {
		return domain.UserToken{}, &InavalidTokenErr
	}

	groups := re.FindStringSubmatch(refreshToken)
	refreshToken = groups[1]

	return handler.service.Refresh(refreshToken)
}

func (handler *AuthRESTHandler) Authenticate(c *gin.Context) (domain.UserToken, error) {
	log.Info().Msg("Start authenticate request")
	log.Info().Msg("Trying to login")
	user, err := handler.Login(c)

	if err == nil {
		log.Info().Msg("Login successful")
		return user, err
	}

	restError, ok := err.(*RestError)
	if ok && restError.Code == UserNotRegistered {
		log.Info().Msg("Trying to register")
		return handler.Register(c)
	}

	log.Error().Err(err).Msg("Other error")
	return domain.UserToken{}, err
}

func (handler *AuthRESTHandler) Me(c *gin.Context) (domain.User, error) {
	userId := c.Request.Header.Get(USER_ID_HEADER)
	if userId == "" {
		return domain.User{}, fmt.Errorf("expected header %q to have exactly one value", USER_ID_HEADER)
	}

	return handler.service.Me(userId)
}

// Utils

func handleError(err error, c *gin.Context) {
	log.Error().Stack().Err(err).Msg("Request error")
	// TODO: Map errors to HTTP status codes
	if rest, ok := err.(*RestError); ok {
		c.JSON(rest.HTTPStatus, gin.H{"error": rest})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
