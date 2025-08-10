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
)

var (
	ErrNoMessage    = errors.New("no message found")
	ErrMessageLimit = errors.New("message limit reached")
)

type Repository interface {
	AddChatroomMessage(context.Context, *Message) error
	GetChatroomMessages(context.Context, string) ([]Message, string, error)
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

func (s *Service) Addclient(c *websocket.Conn, userID, username string) {
	activeUsers := s.getActiveUsers()
	if len(activeUsers) > 0 {
		m := WSMessage{
			Type: typeUserList,
			Data: s.getActiveUsers(),
		}
		data, _ := json.Marshal(m)
		c.WriteMessage(websocket.TextMessage, data)
	}

	u := UserInfo{
		ID:       userID,
		Username: username,
	}
	s.mu.Lock()
	s.clients[c] = &u
	s.mu.Unlock()

	s.broadcast(WSMessage{
		Type: typePresence,
		Data: PresenceMessage{
			Status: statusJoined,
			User:   u,
		},
	})
}

func (s *Service) RemoveClient(c *websocket.Conn, userID, username string) {
	s.mu.Lock()
	delete(s.clients, c)
	s.mu.Unlock()

	u := UserInfo{
		ID:       userID,
		Username: username,
	}

	s.broadcast(WSMessage{
		Type: "presence",
		Data: PresenceMessage{
			Status: statusLeft,
			User:   u,
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
