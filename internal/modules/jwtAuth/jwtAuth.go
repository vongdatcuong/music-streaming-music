package jwtAuth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtService struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	UserID uint64 `json:"user_id"`
}

func NewService(secretKey string, tokenDuration time.Duration) *JwtService {
	return &JwtService{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (s *JwtService) GenerateToken(userID uint64) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenDuration).Unix(),
		},
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s *JwtService) ValidateToken(accessToken string) (*UserClaims, error) {
	var myKey = []byte(s.secretKey)
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Could not validate token")
		}
		return myKey, nil
	})

	if err != nil {
		return &UserClaims{}, fmt.Errorf("could not validate token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok {
		return &UserClaims{}, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *JwtService) InjectToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, "authorization", "Bearer "+accessToken)
}
