// session/session.go
package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

type Store struct {
	RedisClient *redis.Client
}

func NewSessionStore(addr string, password string, db int) *Store {
	return &Store{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (s *Store) SetSession(sessionID string, value string, duration time.Duration) error {
	err := s.RedisClient.Set(sessionID, value, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetSession(cookie *http.Cookie) (string, error) {
	if cookie == nil {
		return "", errors.New("Cookie is nil")
	}
	sessionID := cookie.Value
	val, err := s.RedisClient.Get(sessionID).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (s *Store) DeleteSession(sessionID string) error {
	err := s.RedisClient.Del(sessionID).Err()
	if err != nil {
		return err
	}
	return nil
}

func GenerateSID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.New("Failed to generate SID: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
