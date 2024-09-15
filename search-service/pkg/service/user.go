package service

import (
	"Go_Learn/conf"
	"Go_Learn/pkg/model"
	"Go_Learn/pkg/repo"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTService struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

type UserService struct {
	Repo repo.IRepo
	Conf conf.Config
}

type IUserService interface {
	GenJWTToken(userID int) (string, error)
	ValidateToken(token string) (string, error)
	SearchFile(filter []model.Filter) ([]model.File, error)
}

func NewUserService(repo repo.IRepo, conf conf.Config) IUserService {
	return &UserService{
		Repo: repo,
		Conf: conf,
	}
}

func (s *UserService) GenJWTToken(userID int) (string, error) {
	claims := &JWTService{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "HuyHoang",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.Conf.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *UserService) ValidateToken(token string) (string, error) {
	// validate the claims and verify the signature
	rs, err := jwt.ParseWithClaims(
		token,
		&JWTService{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Conf.SecretKey), nil
		},
	)
	if err != nil {
		return "", err
	}
	// parse the claims from the token
	claims, ok := rs.Claims.(*JWTService)
	if !ok {
		return "", fmt.Errorf("could't parse claims")
	}
	userID := fmt.Sprintf("%v", claims.UserID)
	return userID, nil
}

func (s *UserService) SearchFile(filter []model.Filter) ([]model.File, error) {
	result, err := s.Repo.SearchFiles(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
