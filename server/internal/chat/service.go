package chat

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	newLinePlus = "\n+"
	newLine     = "\n"

	maxMessageLength = 1000
	historyCount     = 20
)

var (
	ErrNoMessage    = errors.New("no message found")
	ErrMessageLimit = errors.New("message limit reached")
)

type Repository interface {
	AddChatroomMessage(context.Context, *Message) error
	GetChatroomMessages(context.Context, string) ([]Message, string, error)
	GetHistory(context.Context, string, int) ([]Message, error)
}

type Service struct {
	repo    Repository
	clients map[*websocket.Conn]*UserInfo
	mu      *sync.RWMutex
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:    repo,
		clients: make(map[*websocket.Conn]*UserInfo),
		mu:      &sync.RWMutex{},
	}
}

func (s *Service) SendChatroomMessage(ctx context.Context, m *Message) error {
	m.Content = strings.TrimSpace(m.Content)
	if m.Content == "" {
		return ErrNoMessage
	}

	if len(m.Content) > maxMessageLength {
		return ErrMessageLimit
	}

	re := regexp.MustCompile(newLinePlus)
	m.Content = re.ReplaceAllString(m.Content, newLine)

	return s.repo.AddChatroomMessage(ctx, m)
}

func (s *Service) broadcast(m WSMessage) {
	data, _ := json.Marshal(m)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for c := range s.clients {
		if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
			c.Close()
			delete(s.clients, c)
		}
	}
}

func (s *Service) Listen(ctx context.Context) {
	lastID := "$"
	for {
		messages, newID, err := s.repo.GetChatroomMessages(ctx, lastID)
		if err != nil {
			log.Printf("chat: redis read error, %v\n", err)
			continue
		}
		lastID = newID
		for _, m := range messages {
			wm := WSMessage{
				Type: "chat",
				Data: m,
			}
			s.broadcast(wm)
		}
	}
}

func (s *Service) Addclient(ctx context.Context, c *websocket.Conn, u *UserInfo) {
	activeUsers := s.getActiveUsers()
	if len(activeUsers) > 0 {
		m := WSMessage{
			Type: typeUserList,
			Data: s.getActiveUsers(),
		}
		data, _ := json.Marshal(m)
		c.WriteMessage(websocket.TextMessage, data)
	}

	s.mu.Lock()
	s.clients[c] = u
	s.mu.Unlock()

	s.broadcast(WSMessage{
		Type: typePresence,
		Data: PresenceMessage{
			Status: statusJoined,
			User:   *u,
		},
	})

	lastID := "+"
	history, _ := s.repo.GetHistory(ctx, lastID, historyCount)
	m := WSMessage{
		Type: typeHistory,
		Data: history,
	}

	data, _ := json.Marshal(m)
	c.WriteMessage(websocket.TextMessage, data)
}

func (s *Service) RemoveClient(c *websocket.Conn, u *UserInfo) {
	s.mu.Lock()
	delete(s.clients, c)
	s.mu.Unlock()

	s.broadcast(WSMessage{
		Type: "presence",
		Data: PresenceMessage{
			Status: statusLeft,
			User:   *u,
		},
	})
}

func (s *Service) getActiveUsers() []*UserInfo {
	var users []*UserInfo

	for _, u := range s.clients {
		users = append(users, u)
	}

	return users
}

func (s *Service) LoadHistoryMessages(ctx context.Context, after string) ([]Message, error) {
	history, err := s.repo.GetHistory(ctx, after, historyCount)
	if err != nil {
		return nil, err
	}

	return history, nil
}
