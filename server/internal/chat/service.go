package chat

import "context"

type Repository interface {
	AddChatroomMessage(context.Context, *Message) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SendChatroomMessage(ctx context.Context, m *Message) error {
	return s.repo.AddChatroomMessage(ctx, m)
}
