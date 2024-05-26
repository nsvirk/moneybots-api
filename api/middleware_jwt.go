package api

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// MWJwtAuthentication middleware authenticates the user using JWT
func MWJwtAuthentication() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(cfg.JWTSigningKey),

		// ErrorHandler is the function to handle error
		// Send custom error message when JWT token is invalid
		ErrorHandler: func(c echo.Context, err error) error {
			return SendError(c, http.StatusUnauthorized, err.Error())
		},

		// SuccessHandler is the function to handle success
		// Set userId and enctoken to the context
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*jwtCustomClaims)
			userId := claims.UserId
			enctoken := claims.Enctoken
			c.Set("user_id", userId)
			c.Set("enctoken", enctoken)
		},
	})
}

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	UserId   string `json:"user_id"`
	Enctoken string `json:"enctoken"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token with custom claims for the user
func GenerateJWT(userId, enctoken string) (string, error) {
	claims := &jwtCustomClaims{
		UserId:   userId,
		Enctoken: enctoken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
			Issuer:    "Moneybots API",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSigningKey))
}
