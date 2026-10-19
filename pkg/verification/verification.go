package verification

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "auth:"

// Store keeps short-lived verification tokens in Redis.
type Store struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewStore(rdb *redis.Client, ttl time.Duration) *Store {
	return &Store{rdb: rdb, ttl: ttl}
}

// TTL reports how long saved tokens stay valid.
func (s *Store) TTL() time.Duration {
	return s.ttl
}

// NewToken returns a URL-safe random token.
func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Save stores a token for a purpose and user. The key combines the public id
// and the token, so a token issued to one user cannot be used for another.
func (s *Store) Save(ctx context.Context, purpose, publicID, token string) error {
	return s.rdb.Set(ctx, key(purpose, publicID, token), "1", s.ttl).Err()
}

// Take atomically consumes the token, reporting whether it was valid.
func (s *Store) Take(ctx context.Context, purpose, publicID, token string) (bool, error) {
	_, err := s.rdb.GetDel(ctx, key(purpose, publicID, token)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func key(purpose, publicID, token string) string {
	return fmt.Sprintf("%s%s:%s:%s", keyPrefix, purpose, publicID, token)
}
