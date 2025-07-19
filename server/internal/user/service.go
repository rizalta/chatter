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

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Repository interface {
	CreateUser(u *User) error
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
		return nil, fmt.Errorf("user: failed to hash password, %v", err)
	}

	u := User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := s.repo.CreateUser(&u); err != nil {
		return nil, fmt.Errorf("user: failed to create user, %v", err)
	}

	return &u, err
}

func (s *Service) Login(username, password string) (string, error) {
	u, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("user: failed to check password, %v", err)
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
