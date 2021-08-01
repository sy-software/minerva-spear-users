package service

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	mvdatetime "github.com/sy-software/minerva-go-utils/datetime"
	"github.com/sy-software/minerva-spear-users/internal/core/domain"
	"github.com/sy-software/minerva-spear-users/mocks"
)

func TestUserRegistration(t *testing.T) {
	config := domain.DefaultConfig()
	config.Token.PrivateKey = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAqKIQvqxSMVaplK3aBdRJyRYf5IhvCPH4IF2DlmmPd6ilFzH3\njJTjnDt2+6GmAQGhHau6LPExdSHmbLCck124JS8mbis83zAOQ3hfqmvgvAO4smAE\n3dxE4XX2SRsFl1aOV6oPM4ckZgTLMDNkocxfo4TVV4Yg3ycf74MKfh+47hwJAyLZ\nJFImCnheDij2YARsEHSKAdX9iEP9IFSDqhX0+XJyGxB07wHWX6fHjcXaUfKq68CI\nrX5d22m8ZN1zCmJYNTfJNpfvquy+uosSUDNU4W9WFHJmOJS6jE5lQYcbCCROlZRW\njuBfy9UJnl4jYDgWjClYYU3qrv1UGD05Vn9eIwIDAQABAoIBAE7aQYQ3ZdOmV3Or\ne5hgNQRvcQhW97yyELlpoN9Tiv+D/3aCKeQ1ttzWPYPaiZpM3b7XDx52xg6khG/s\ngbqzByl0C79WPoeKnBDWl71D5nlkMBhQp9XqatcWZsy2cv3aPoMlhSguGEoQEcb/\nMR4rR8lZkrzzfil6zQcdOmnRgZLtDL2l2gC1NIZtU/4bS3fWe7AxcnT7a6UIrxyL\nkfzOuxgUawLR5Bth/mCDbmJgsDXndT4WNO3CWBzEZ0WLXdkyAbFgPzIPMoeR+Y68\nYLBw1OhrzrMP4FrNYs6cua3HGR2PzcH3/qec1i8grMyVwCUWVLcRn36ZGZnrhjIy\nKO9cCgkCgYEAzA8QKEXAVwAyxrjPC+cwywBLk7BsYn+4MKUw3z9h8DgzKLIojxCC\nw4WSO+mZ+VFOS9lOkGPG5/d0s/8ZWo8WrOApa6pXiz5JC8IeOEh3xPNN9xpJPW7P\nbjcbaQ+gppj/fTullfrEo4O3/5YDldy7enk8xDrhtbC9iS25QNDm6s8CgYEA046U\nIGDvpOx3iv3G1dqXPFQGSe1JJ4wytqAwz8chlj63R2lIV6LiIlgEqdd0TyZag8yB\nCBGE0j+TmDc6+4xHqRCwkKL+1e2s3tR0X/Vs4CP9PXiywAbY0vhI2CtguIuBYf1Z\nruFZ5bMkf+sT4evJlSMRrmdtKoAumYXXfYJFXG0CgYEAkQb7osPAKZU4gUgDzx/m\n68Av9q1iuravP9OH4oL3pnUq1veYH+XKKhAamH40Mp/4l6vATJq9WUvkI7FgYZ5k\nrUU76wtL4OjJnZO/Sp0mklGhzcde2kyRHHIKBydWNFF085qa2vc5HkWVVg9WSQJy\nNF9KMuTuWeVdL8vRaCGQnL0CgYEAm7OyHWp6td071lYUw0xQRpxozHwRfUPYB0U6\n55FdjOC3r50zGxzMZg510DK8bYyCzcHzrWaHZN5Z2Iu9o2mJTEr2SF1ORVDaDF49\nEGrnKMgUF+v/Uwk3B36ozkCOvQQfw2jdWrKMoVwJnwP67CnHgTYAS2XfmIoiwecZ\nxEvelLkCgYEAgGqcnf4RRoIU7joEjXCorBDebDubj5C0ektbenu9rNTiCSiU7T+b\np6S9CGtMfMbl3b/w2JlojbElGvkqaBdkbWhsRpRzlRh8IC0/Gn4l03i2u9YNhVia\nQrdb/UxBM3vRVzCf216QcgCHGNQXKKQtBHu71+cFMUp5sExm6XRQ0Qc=\n-----END RSA PRIVATE KEY-----"
	config.Token.PublicKey = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqKIQvqxSMVaplK3aBdRJ\nyRYf5IhvCPH4IF2DlmmPd6ilFzH3jJTjnDt2+6GmAQGhHau6LPExdSHmbLCck124\nJS8mbis83zAOQ3hfqmvgvAO4smAE3dxE4XX2SRsFl1aOV6oPM4ckZgTLMDNkocxf\no4TVV4Yg3ycf74MKfh+47hwJAyLZJFImCnheDij2YARsEHSKAdX9iEP9IFSDqhX0\n+XJyGxB07wHWX6fHjcXaUfKq68CIrX5d22m8ZN1zCmJYNTfJNpfvquy+uosSUDNU\n4W9WFHJmOJS6jE5lQYcbCCROlZRWjuBfy9UJnl4jYDgWjClYYU3qrv1UGD05Vn9e\nIwIDAQAB\n-----END PUBLIC KEY-----"
	called := false

	registerReq := domain.Register{
		Name:     "Tony Stark",
		Username: "IronMan",
		Picture:  "https://picture.com/ironman",
		Role:     "hero",
		Provider: "StarkIndustries",
		TokenID:  "tokenId",
	}

	expectedInfo := domain.User{
		Id:       "newid",
		Name:     "Tony Stark",
		Username: "IronMan",
		Picture:  "https://picture.com/ironman",
	}

	repo := mocks.UserRepo{
		CreateInterceptor: func(user domain.Register) (domain.User, error) {
			called = true

			if !cmp.Equal(user, registerReq) {
				t.Errorf("Expected create to be called with: %+v Got: %+v", registerReq, user)
			}
			return expectedInfo, nil
		},
	}

	service := NewAuthService(&repo, config)
	now := mvdatetime.UnixUTCNow()
	expire := now.Add(time.Duration(config.Token.Duration) * time.Second)
	refreshExpire := now.Add(time.Duration(config.Token.RefreshDuration) * time.Second)
	token, err := service.Register(registerReq)

	if !called {
		t.Error("Expected repo.Create to be called")
	}

	if err != nil {
		t.Errorf("Expected register without error, got: %v", err)
	}

	k, err := service.config.Token.KeyPair()
	decoded, err := jwt.Parse([]byte(token.AccessToken), jwt.WithVerify(jwa.RS256, k.PublicKey), jwt.WithValidate(true))

	if err != nil {
		t.Errorf("Expected JWT to be decoded without error, got: %v", err)
	}

	if !decoded.Expiration().Equal(token.ExpireTime) {
		t.Errorf("Decoded token expire date should be: %v got: %v", token.ExpireTime, decoded.Expiration())
	}

	use, ok := decoded.Get("use")
	if !ok || use != "access" {
		t.Errorf("Expected token use to be access got: %q", use)
	}

	decodedRefresh, err := jwt.Parse([]byte(token.RefreshToken), jwt.WithVerify(jwa.RS256, k.PublicKey), jwt.WithValidate(true))

	if err != nil {
		t.Errorf("Expected JWT to be decoded without error, got: %v", err)
	}

	if decodedRefresh.Expiration().Before(refreshExpire) {
		t.Errorf("Decoded refresh token expire date should be at least: %v got: %v", refreshExpire, decodedRefresh.Expiration())
	}

	use, ok = decodedRefresh.Get("use")
	if !ok || use != "refresh" {
		t.Errorf("Expected token use to be refresh got: %q", use)
	}

	if token.TokenType != "Bearer" {
		t.Errorf("Expected token type Bearer got: %q", token.TokenType)
	}

	if token.ExpireTime.Before(expire) {
		t.Errorf("Expected expire time to be at least: %v got: %v", expire, token.ExpireTime)
	}

	if !cmp.Equal(expectedInfo, token.Info) {
		t.Errorf("Expected user info to be: %+v; got: %+v", expectedInfo, token.Info)
	}
}
