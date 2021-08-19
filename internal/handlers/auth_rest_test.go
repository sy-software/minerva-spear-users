package handlers

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sy-software/minerva-spear-users/internal/core/domain"
	"github.com/sy-software/minerva-spear-users/internal/core/service"
	"github.com/sy-software/minerva-spear-users/mocks"
)

const PRIVATE_KEY = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAqKIQvqxSMVaplK3aBdRJyRYf5IhvCPH4IF2DlmmPd6ilFzH3\njJTjnDt2+6GmAQGhHau6LPExdSHmbLCck124JS8mbis83zAOQ3hfqmvgvAO4smAE\n3dxE4XX2SRsFl1aOV6oPM4ckZgTLMDNkocxfo4TVV4Yg3ycf74MKfh+47hwJAyLZ\nJFImCnheDij2YARsEHSKAdX9iEP9IFSDqhX0+XJyGxB07wHWX6fHjcXaUfKq68CI\nrX5d22m8ZN1zCmJYNTfJNpfvquy+uosSUDNU4W9WFHJmOJS6jE5lQYcbCCROlZRW\njuBfy9UJnl4jYDgWjClYYU3qrv1UGD05Vn9eIwIDAQABAoIBAE7aQYQ3ZdOmV3Or\ne5hgNQRvcQhW97yyELlpoN9Tiv+D/3aCKeQ1ttzWPYPaiZpM3b7XDx52xg6khG/s\ngbqzByl0C79WPoeKnBDWl71D5nlkMBhQp9XqatcWZsy2cv3aPoMlhSguGEoQEcb/\nMR4rR8lZkrzzfil6zQcdOmnRgZLtDL2l2gC1NIZtU/4bS3fWe7AxcnT7a6UIrxyL\nkfzOuxgUawLR5Bth/mCDbmJgsDXndT4WNO3CWBzEZ0WLXdkyAbFgPzIPMoeR+Y68\nYLBw1OhrzrMP4FrNYs6cua3HGR2PzcH3/qec1i8grMyVwCUWVLcRn36ZGZnrhjIy\nKO9cCgkCgYEAzA8QKEXAVwAyxrjPC+cwywBLk7BsYn+4MKUw3z9h8DgzKLIojxCC\nw4WSO+mZ+VFOS9lOkGPG5/d0s/8ZWo8WrOApa6pXiz5JC8IeOEh3xPNN9xpJPW7P\nbjcbaQ+gppj/fTullfrEo4O3/5YDldy7enk8xDrhtbC9iS25QNDm6s8CgYEA046U\nIGDvpOx3iv3G1dqXPFQGSe1JJ4wytqAwz8chlj63R2lIV6LiIlgEqdd0TyZag8yB\nCBGE0j+TmDc6+4xHqRCwkKL+1e2s3tR0X/Vs4CP9PXiywAbY0vhI2CtguIuBYf1Z\nruFZ5bMkf+sT4evJlSMRrmdtKoAumYXXfYJFXG0CgYEAkQb7osPAKZU4gUgDzx/m\n68Av9q1iuravP9OH4oL3pnUq1veYH+XKKhAamH40Mp/4l6vATJq9WUvkI7FgYZ5k\nrUU76wtL4OjJnZO/Sp0mklGhzcde2kyRHHIKBydWNFF085qa2vc5HkWVVg9WSQJy\nNF9KMuTuWeVdL8vRaCGQnL0CgYEAm7OyHWp6td071lYUw0xQRpxozHwRfUPYB0U6\n55FdjOC3r50zGxzMZg510DK8bYyCzcHzrWaHZN5Z2Iu9o2mJTEr2SF1ORVDaDF49\nEGrnKMgUF+v/Uwk3B36ozkCOvQQfw2jdWrKMoVwJnwP67CnHgTYAS2XfmIoiwecZ\nxEvelLkCgYEAgGqcnf4RRoIU7joEjXCorBDebDubj5C0ektbenu9rNTiCSiU7T+b\np6S9CGtMfMbl3b/w2JlojbElGvkqaBdkbWhsRpRzlRh8IC0/Gn4l03i2u9YNhVia\nQrdb/UxBM3vRVzCf216QcgCHGNQXKKQtBHu71+cFMUp5sExm6XRQ0Qc=\n-----END RSA PRIVATE KEY-----"
const PUBLIC_KEY = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqKIQvqxSMVaplK3aBdRJ\nyRYf5IhvCPH4IF2DlmmPd6ilFzH3jJTjnDt2+6GmAQGhHau6LPExdSHmbLCck124\nJS8mbis83zAOQ3hfqmvgvAO4smAE3dxE4XX2SRsFl1aOV6oPM4ckZgTLMDNkocxf\no4TVV4Yg3ycf74MKfh+47hwJAyLZJFImCnheDij2YARsEHSKAdX9iEP9IFSDqhX0\n+XJyGxB07wHWX6fHjcXaUfKq68CIrX5d22m8ZN1zCmJYNTfJNpfvquy+uosSUDNU\n4W9WFHJmOJS6jE5lQYcbCCROlZRWjuBfy9UJnl4jYDgWjClYYU3qrv1UGD05Vn9e\nIwIDAQAB\n-----END PUBLIC KEY-----"

func TestLoginEndpoint(t *testing.T) {
	config := domain.DefaultConfig()
	config.Token.PrivateKey = PRIVATE_KEY
	config.Token.PublicKey = PUBLIC_KEY

	repo := mocks.UserRepo{
		GetByUsernameInterceptor: func(username string) (domain.User, error) {
			return domain.User{
				Id:       "newid",
				Name:     "Tony Stark",
				Username: "IronMan",
				Picture:  "https://picture.com/ironman",
			}, nil
		},
	}

	service := service.NewAuthService(&repo, config)

	handler := NewAuthRESTHandler(&config, service)

	userInfo := `
	{
		"username": "IronMan",
		"provder": "StarkIndustries",
		"tokenID": "myTokenId"
	}
	`
	headers := http.Header{}
	headers.Add(USER_INFO_HEADER, base64.StdEncoding.EncodeToString([]byte(userInfo)))
	context := gin.Context{
		Request: &http.Request{
			Header: headers,
		},
	}

	token, err := handler.Login(&context)

	if err != nil {
		t.Errorf("Expected to login without error, got: %v", err)
	}

	if token.Info.Username != "IronMan" {
		t.Errorf("Expected token for username: IronMan got: %q", token.Info.Username)
	}
}

func TestRegisterEndpoint(t *testing.T) {
	config := domain.DefaultConfig()
	config.Token.PrivateKey = PRIVATE_KEY
	config.Token.PublicKey = PUBLIC_KEY

	repo := mocks.UserRepo{
		CreateInterceptor: func(user domain.Register) (domain.User, error) {
			return domain.User{
				Id:       "newid",
				Name:     "Tony Stark",
				Username: "IronMan",
				Picture:  "https://picture.com/ironman",
			}, nil
		},
	}

	service := service.NewAuthService(&repo, config)

	handler := NewAuthRESTHandler(&config, service)

	userInfo := `
	{
		"username": "IronMan",
		"name": "Tony Stark",
		"picture": "https://picture.com/tony",
		"role": "hero",
		"provder": "StarkIndustries",
		"tokenID": "myTokenId"
	}
	`
	headers := http.Header{}
	headers.Add(USER_INFO_HEADER, base64.StdEncoding.EncodeToString([]byte(userInfo)))
	context := gin.Context{
		Request: &http.Request{
			Header: headers,
		},
	}

	token, err := handler.Register(&context)

	if err != nil {
		t.Errorf("Expected to register without error, got: %v", err)
	}

	if token.Info.Username != "IronMan" {
		t.Errorf("Expected token for username: IronMan got: %q", token.Info.Username)
	}
}

func TestRefreshEndpoint(t *testing.T) {
	config := domain.DefaultConfig()
	config.Token.PrivateKey = PRIVATE_KEY
	config.Token.PublicKey = PUBLIC_KEY

	repo := mocks.UserRepo{
		GetByUsernameInterceptor: func(username string) (domain.User, error) {
			return domain.User{
				Id:       "newid",
				Name:     "Tony Stark",
				Username: "IronMan",
				Picture:  "https://picture.com/ironman",
			}, nil
		},
		GetByIdInterceptor: func(id string) (domain.User, error) {
			return domain.User{
				Id:       "newid",
				Name:     "Tony Stark",
				Username: "IronMan",
				Picture:  "https://picture.com/ironman",
			}, nil
		},
	}

	service := service.NewAuthService(&repo, config)

	handler := NewAuthRESTHandler(&config, service)

	userInfo := `
	{
		"username": "IronMan",
		"provder": "StarkIndustries",
		"tokenID": "myTokenId"
	}
	`

	headers := http.Header{}
	headers.Add(USER_INFO_HEADER, base64.StdEncoding.EncodeToString([]byte(userInfo)))
	context := gin.Context{
		Request: &http.Request{
			Header: headers,
		},
	}

	login, _ := handler.Login(&context)
	// Wait a few seconds to have a token with different expiration
	time.Sleep(1 * time.Second)

	headers = http.Header{}
	headers.Add("Authorization", "Bearer "+login.RefreshToken)
	context = gin.Context{
		Request: &http.Request{
			Header: headers,
		},
	}
	token, err := handler.Refresh(&context)

	if err != nil {
		t.Errorf("Expected to refresh token without error, got: %v", err)
	}

	if token.Info.Username != "IronMan" {
		t.Errorf("Expected token for username: IronMan got: %q", token.Info.Username)
	}

	if login.AccessToken == token.AccessToken {
		t.Errorf("Expected a new access token got the same as login")
	}

	if login.RefreshToken == token.RefreshToken {
		t.Errorf("Expected a new refresh access token got the same as login")
	}
}

func TestAuthenticateEndpoint(t *testing.T) {
	t.Run("Test login", func(t *testing.T) {
		config := domain.DefaultConfig()
		config.Token.PrivateKey = PRIVATE_KEY
		config.Token.PublicKey = PUBLIC_KEY

		repo := mocks.UserRepo{
			GetByUsernameInterceptor: func(username string) (domain.User, error) {
				return domain.User{
					Id:       "newid",
					Name:     "Tony Stark",
					Username: "IronMan",
					Picture:  "https://picture.com/ironman",
				}, nil
			},
		}

		service := service.NewAuthService(&repo, config)

		handler := NewAuthRESTHandler(&config, service)

		userInfo := `
			{
				"username": "IronMan",
				"provder": "StarkIndustries",
				"tokenID": "myTokenId"
			}
		`
		headers := http.Header{}
		headers.Add(USER_INFO_HEADER, base64.StdEncoding.EncodeToString([]byte(userInfo)))
		context := gin.Context{
			Request: &http.Request{
				Header: headers,
			},
		}

		token, err := handler.Authenticate(&context)

		if err != nil {
			t.Errorf("Expected to login without error, got: %v", err)
		}

		if token.Info.Username != "IronMan" {
			t.Errorf("Expected token for username: IronMan got: %q", token.Info.Username)
		}
	})

	t.Run("Test register", func(t *testing.T) {
		config := domain.DefaultConfig()
		config.Token.PrivateKey = PRIVATE_KEY
		config.Token.PublicKey = PUBLIC_KEY

		repo := mocks.UserRepo{
			GetByUsernameInterceptor: func(username string) (domain.User, error) {
				return domain.User{}, errors.New("not_found")
			},
			CreateInterceptor: func(user domain.Register) (domain.User, error) {
				return domain.User{
					Id:       "newid",
					Name:     "Tony Stark",
					Username: "IronMan",
					Picture:  "https://picture.com/ironman",
				}, nil
			},
		}

		service := service.NewAuthService(&repo, config)

		handler := NewAuthRESTHandler(&config, service)

		userInfo := `
		{
			"username": "IronMan",
			"name": "Tony Stark",
			"picture": "https://picture.com/tony",
			"role": "hero",
			"provder": "StarkIndustries",
			"tokenID": "myTokenId"
		}
		`

		headers := http.Header{}
		headers.Add(USER_INFO_HEADER, base64.StdEncoding.EncodeToString([]byte(userInfo)))
		context := gin.Context{
			Request: &http.Request{
				Header: headers,
			},
		}

		token, err := handler.Authenticate(&context)

		if err != nil {
			t.Errorf("Expected to register without error, got: %v", err)
		}

		if token.Info.Username != "IronMan" {
			t.Errorf("Expected token for username: IronMan got: %q", token.Info.Username)
		}
	})
}

func TestMeEndpoint(t *testing.T) {
	config := domain.DefaultConfig()
	config.Token.PrivateKey = PRIVATE_KEY
	config.Token.PublicKey = PUBLIC_KEY

	called := false
	repo := mocks.UserRepo{
		GetByUsernameInterceptor: func(username string) (domain.User, error) {
			return domain.User{
				Id:       "newid",
				Name:     "Tony Stark",
				Username: "IronMan",
				Picture:  "https://picture.com/ironman",
			}, nil
		},
		GetByIdInterceptor: func(id string) (domain.User, error) {
			called = true

			if id != "newid" {
				t.Errorf("Expected id to be \"newid\" got: %q", id)
			}

			return domain.User{
				Id:       "newid",
				Name:     "Tony Stark",
				Username: "IronMan",
				Picture:  "https://picture.com/ironman",
			}, nil
		},
	}

	service := service.NewAuthService(&repo, config)

	handler := NewAuthRESTHandler(&config, service)

	body := `
	{
		"username": "IronMan",
		"provder": "StarkIndustries",
		"tokenID": "myTokenId"
	}
	`
	bodyReader := io.NopCloser(strings.NewReader(body))
	context := gin.Context{
		Request: &http.Request{
			Header: http.Header{},
			Body:   bodyReader,
		},
	}

	login, _ := handler.Login(&context)
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+login.AccessToken)
	headers.Add(USER_ID_HEADER, "newid")
	context = gin.Context{
		Request: &http.Request{
			Header: headers,
		},
	}
	info, err := handler.Me(&context)

	if err != nil {
		t.Errorf("Expected to /me to return without error, got: %v", err)
	}

	if !called {
		t.Error("Expected repo.GetById to be called")
	}

	if info.Username != "IronMan" {
		t.Errorf("Expected token for username: IronMan got: %q", info.Username)
	}
}

func TestErrors(t *testing.T) {
	config := domain.DefaultConfig()
	config.Token.PrivateKey = PRIVATE_KEY
	config.Token.PublicKey = PUBLIC_KEY

	t.Run("Test invalid token error", func(t *testing.T) {
		repo := mocks.UserRepo{
			GetByUsernameInterceptor: func(username string) (domain.User, error) {
				return domain.User{
					Id:       "newid",
					Name:     "Tony Stark",
					Username: "IronMan",
					Picture:  "https://picture.com/ironman",
				}, nil
			},
		}

		service := service.NewAuthService(&repo, config)

		handler := NewAuthRESTHandler(&config, service)
		headers := http.Header{}
		headers.Add("Authorization", "Not A Bearer token")
		context := gin.Context{
			Request: &http.Request{
				Header: headers,
			},
		}

		_, err := handler.Refresh(&context)

		if err == nil {
			t.Errorf("Expected an error got nil")
		}

		parsed, ok := err.(*RestError)

		if !ok {
			t.Errorf("Expected error of type RestError got: %v", err)
		}

		if parsed.Code != InvalidToken {
			t.Errorf("Expected error code: %d got: %d", InvalidToken, parsed.Code)
		}
	})

	t.Run("Test user not registered error", func(t *testing.T) {
		repo := mocks.UserRepo{
			GetByUsernameInterceptor: func(username string) (domain.User, error) {
				return domain.User{}, errors.New("not_found")
			},
		}

		service := service.NewAuthService(&repo, config)

		handler := NewAuthRESTHandler(&config, service)

		userInfo := `
		{
			"username": "IronMan",
			"provder": "StarkIndustries",
			"tokenID": "myTokenId"
		}
		`
		headers := http.Header{}
		headers.Add(USER_INFO_HEADER, base64.StdEncoding.EncodeToString([]byte(userInfo)))
		context := gin.Context{
			Request: &http.Request{
				Header: headers,
			},
		}

		_, err := handler.Login(&context)

		if err == nil {
			t.Errorf("Expected an error got nil")
		}

		parsed, ok := err.(*RestError)

		if !ok {
			t.Errorf("Expected error of type RestError got: %v", err)
		}

		if parsed.Code != UserNotRegistered {
			t.Errorf("Expected error code: %d got: %d", UserNotRegistered, parsed.Code)
		}
	})

	t.Run("Test user already registered error", func(t *testing.T) {
		repo := mocks.UserRepo{
			CreateInterceptor: func(user domain.Register) (domain.User, error) {
				return domain.User{}, errors.New("duplicated_value")
			},
		}

		service := service.NewAuthService(&repo, config)

		handler := NewAuthRESTHandler(&config, service)

		userInfo := `
		{
			"username": "IronMan",
			"name": "Tony Stark",
			"picture": "https://picture.com/tony",
			"role": "hero",
			"provder": "StarkIndustries",
			"tokenID": "myTokenId"
		}
		`

		headers := http.Header{}
		headers.Add(USER_INFO_HEADER, base64.StdEncoding.EncodeToString([]byte(userInfo)))
		context := gin.Context{
			Request: &http.Request{
				Header: headers,
			},
		}

		_, err := handler.Register(&context)

		if err == nil {
			t.Errorf("Expected an error got nil")
		}

		parsed, ok := err.(*RestError)

		if !ok {
			t.Errorf("Expected error of type RestError got: %v", err)
		}

		if parsed.Code != UserAlreadyRegistered {
			t.Errorf("Expected error code: %d got: %d", UserAlreadyRegistered, parsed.Code)
		}
	})

	t.Run("Test invalid user info header", func(t *testing.T) {
		repo := mocks.UserRepo{
			CreateInterceptor: func(user domain.Register) (domain.User, error) {
				return domain.User{}, errors.New("duplicated_value")
			},
		}

		service := service.NewAuthService(&repo, config)

		handler := NewAuthRESTHandler(&config, service)

		userInfo := "not json"

		headers := http.Header{}
		headers.Add(USER_INFO_HEADER, base64.StdEncoding.EncodeToString([]byte(userInfo)))
		context := gin.Context{
			Request: &http.Request{
				Header: headers,
			},
		}

		_, err := handler.Register(&context)

		if err == nil {
			t.Errorf("Expected an error got nil")
		}

		parsed, ok := err.(*RestError)

		if !ok {
			t.Errorf("Expected error of type RestError got: %v", err)
		}

		if parsed.Code != InavalidRequest {
			t.Errorf("Expected error code: %d got: %d", InavalidRequest, parsed.Code)
		}

		_, err = handler.Login(&context)

		if err == nil {
			t.Errorf("Expected an error got nil")
		}

		parsed, ok = err.(*RestError)

		if !ok {
			t.Errorf("Expected error of type RestError got: %v", err)
		}

		if parsed.Code != InavalidRequest {
			t.Errorf("Expected error code: %d got: %d", InavalidRequest, parsed.Code)
		}
	})
}
