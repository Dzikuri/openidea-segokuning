package helper

import (
	"os"
	"strconv"
	"time"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
	jwt "github.com/golang-jwt/jwt/v5"
)

func JwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	return secret
}

func JwtExpireAt() int {
	expireAt := os.Getenv("JWT_EXPIRE_AT")
	if expireAt == "" {
		expireAt = "2"
	}

	expire, _ := strconv.Atoi(expireAt)

	return expire
}

type JwtCustomClaims struct {
	Name string `json:"name"`
	Id   string `json:"userId"`
	jwt.RegisteredClaims
}

func JwtGenerateToken(request *model.UserResponse) (string, error) {

	// Generate Claims object
	jwtClaims := JwtCustomClaims{
		Name: request.Name,
		Id:   request.Id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}

	// Create token with claims
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(JwtSecret()))
	if err != nil {
		return "", err
	}

	return t, err
}

func VerifyJwt(tokenString string, claims jwt.Claims, secret string) error {
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(JwtSecret()), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return model.ErrUnauthorize
		}
		return err

	}
	if !tkn.Valid {
		return model.ErrUnauthorize
	}

	return nil
}
