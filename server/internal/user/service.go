package user

import (
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(username, password string) (*User, error)
	GetUserByUsername(username string) (*User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) Register(username, password string) (*User, error) {
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

	// TODO genrate jwt
	return "jwt", nil
}
