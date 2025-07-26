package user

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const jwtExpirationTime = 24 * time.Hour

var (
	ErrUsernameLength        = errors.New("username must be between 4 and 20 characters")
	ErrUsernameStart         = errors.New("username must start with a letter")
	ErrUsernameContains      = errors.New("username can only contain letters, numbers, and underscores")
	ErrUsernameAlreadyExists = errors.New("username already exists")

	ErrPasswordLength    = errors.New("password must be at least 8 characters long")
	ErrPasswordUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordLowercase = errors.New("password must contain at least one lowercase letter")
	ErrPasswordDigit     = errors.New("password must contain at least one digit")
	ErrPasswordSpecial   = errors.New("password must contain at least one special character")

	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Repository interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
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

func (s *Service) Register(ctx context.Context, username, password string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	_, err := s.repo.GetUserByUsername(ctx, username)
	if err == nil {
		return ErrUsernameAlreadyExists
	}

	if err := validatePassword(password); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("user: failed to hash password, %v", err)
	}

	u := User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := s.repo.CreateUser(ctx, &u); err != nil {
		return fmt.Errorf("user: failed to create user, %v", err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	u, err := s.repo.GetUserByUsername(ctx, username)
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

func validateUsername(username string) error {
	if len(username) < 4 || len(username) > 20 {
		return ErrUsernameLength
	}

	if matched, _ := regexp.MatchString(`^[a-zA-Z]`, username); !matched {
		return ErrUsernameStart
	}

	if matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, username); !matched {
		return ErrUsernameContains
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordLength
	}

	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return ErrPasswordUppercase
	}

	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return ErrPasswordLowercase
	}

	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return ErrPasswordDigit
	}

	if matched, _ := regexp.MatchString(`[\W_]`, password); !matched {
		return ErrPasswordSpecial
	}

	return nil
}
