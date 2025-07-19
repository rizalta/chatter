package user

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const jwtExpirationTime = 24 * time.Hour

var ErrUserAlreadyExists = errors.New("user already exists")

type Repository interface {
	CreateUser(username, password string) (*User, error)
	GetUserByUsername(username string) (*User, error)
}

type Service struct {
	repo       Repository
	privateKey ed25519.PrivateKey
}

func NewService(repo Repository, privateKey ed25519.PrivateKey) *Service {
	return &Service{
		repo:       repo,
		privateKey: privateKey,
	}
}

func (s *Service) Register(username, password string) (*User, error) {
	_, err := s.repo.GetUserByUsername(username)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateUser(username, string(hashedPassword))
}

func (s *Service) Login(username, password string) (string, error) {
	u, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", err
	}

	jwtToken, err := s.generateToken(u.ID)
	if err != nil {
		return "", fmt.Errorf("user: error genrating jwt token: %v", err)
	}

	return jwtToken, nil
}

func (s *Service) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	return token.SignedString(s.privateKey)
}
