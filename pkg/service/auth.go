package service

import (
	"banner"
	"banner/pkg/repository"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	signingKey = "adfa6464aE"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user banner.User) (int, error) {
	return s.repo.CreateUser(user)
}

func (s *AuthService) CheckNickNameAndEmail(nickname, email string) (int, error) {
	return s.repo.CheckNickNameAndEmail(nickname, email)
}

func (s *AuthService) GetPasswordHash(nickname string) (string, error) {
	return s.repo.GetPasswordHash(nickname)
}

func (s *AuthService) GenerateToken(nickname, passwordHash string) (string, error) {
	user, err := s.repo.GetUser(nickname, passwordHash)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.Role,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(accessToken *jwt.Token) (interface{}, error) {
		if _, ok := accessToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, claims.Role, nil
}
