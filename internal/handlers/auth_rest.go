package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sy-software/minerva-spear-users/internal/core/domain"
	"github.com/sy-software/minerva-spear-users/internal/core/ports"
)

const BEARER_REGEX = "^Bearer (.+)$"

const USER_ID_HEADER = "X-USER-ID"

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
	router.POST(fmt.Sprintf("%s/login", handler.config.APIPrefix), func(c *gin.Context) {
		token, err := handler.Login(c)

		if err != nil {
			handleError(err, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": token})
	})

	router.POST(fmt.Sprintf("%s/refresh", handler.config.APIPrefix), func(c *gin.Context) {
		token, err := handler.Refresh(c)

		if err != nil {
			handleError(err, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": token})
	})

	router.POST(fmt.Sprintf("%s/register", handler.config.APIPrefix), func(c *gin.Context) {
		token, err := handler.Register(c)

		if err != nil {
			handleError(err, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": token})
	})

	router.GET(fmt.Sprintf("%s/me", handler.config.APIPrefix), func(c *gin.Context) {
		user, err := handler.Me(c)

		if err != nil {
			handleError(err, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": user})
	})
}

func (handler *AuthRESTHandler) Login(c *gin.Context) (domain.UserToken, error) {
	var body domain.Login
	err := c.BindJSON(&body)

	if err != nil {
		return domain.UserToken{}, &InavalidBodyErr
	}

	user, err := handler.service.Login(body)

	if err != nil {
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
	var body domain.Register

	err := c.BindJSON(&body)

	if err != nil {
		return domain.UserToken{}, &InavalidBodyErr
	}

	user, err := handler.service.Register(body)

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

func (handler *AuthRESTHandler) Me(c *gin.Context) (domain.User, error) {
	if len(c.Request.Header[USER_ID_HEADER]) != 1 {
		return domain.User{}, fmt.Errorf("expected header %q to have exactly one value", USER_ID_HEADER)
	}

	userId := c.Request.Header[USER_ID_HEADER][0]

	return handler.service.Me(userId)
}

// Utils

func handleError(err error, c *gin.Context) {
	log.Error().Stack().Err(err).Msg("Request error:")
	// TODO: Map errors to HTTP status codes
	if rest, ok := err.(*RestError); ok {
		c.JSON(rest.HTTPStatus, gin.H{"error": rest})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
