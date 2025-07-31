package database

import (
	"chatter/server/internal/chat"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ChatRepo struct {
	db *redis.Client
}

func NewChatRepo(db *redis.Client) *ChatRepo {
	return &ChatRepo{db: db}
}

const chatroomKey = "chatroom"

func (r *ChatRepo) AddChatroomMessage(ctx context.Context, m *chat.Message) error {
	m.Timestamp = time.Now().UTC()
	_, err := r.db.XAdd(ctx, &redis.XAddArgs{
		Stream: chatroomKey,
		Values: messageToMap(m),
	}).Result()

	return err
}

func (r *ChatRepo) GetChatroomMessages(ctx context.Context, after string) ([]chat.Message, string, error) {
	streams, err := r.db.XRead(ctx, &redis.XReadArgs{
		Streams: []string{chatroomKey, after},
		Count:   1,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, after, err
	}

	messages := streamsToMessages(streams)
	lastID := messages[len(messages)-1].ID
	fmt.Println(messages)

	return streamsToMessages(streams), lastID, nil
}

func (r *ChatRepo) AddPrivateMessage(ctx context.Context, m *chat.Message) error {
	m.Timestamp = time.Now().UTC()
	key := sortedKey(m.From, m.To)
	_, err := r.db.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		Values: messageToMap(m),
	}).Result()

	return err
}

func sortedKey(user1, user2 string) string {
	if user1 < user2 {
		return fmt.Sprintf("%s:%s", user1, user2)
	}
	return fmt.Sprintf("%s:%s", user2, user1)
}

func messageToMap(m *chat.Message) map[string]string {
	return map[string]string{
		"from":      m.From,
		"to":        m.To,
		"content":   m.Content,
		"timestamp": m.Timestamp.Format(time.RFC3339),
	}
}

func streamsToMessages(streams []redis.XStream) []chat.Message {
	var messages []chat.Message

	for _, stream := range streams {
		for _, entry := range stream.Messages {
			m := chat.Message{
				ID:      entry.ID,
				From:    entry.Values["from"].(string),
				Content: entry.Values["content"].(string),
			}

			if to, ok := entry.Values["to"].(string); ok {
				m.To = to
			}

			if tsStr, ok := entry.Values["timestamp"].(string); ok {
				if ts, err := time.Parse(time.RFC3339, tsStr); err == nil {
					m.Timestamp = ts
				}
			}

			messages = append(messages, m)
		}
	}

	return messages
}
