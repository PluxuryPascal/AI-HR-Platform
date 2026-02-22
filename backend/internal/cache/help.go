package cache

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func fullKey[T any](m *Manager, k Key[T], id string) string {
	if m.prefix != "" {
		return fmt.Sprintf("%s:%s:%s", m.prefix, k.namespace, id)
	}

	return fmt.Sprintf("%s:%s", k.namespace, id)
}

func encode[T any](v T) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("cache: marshal: %w", err)
	}

	return b, nil
}

func decode[T any](raw string) (T, error) {
	var v T

	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return v, fmt.Errorf("cache: unmarshal: %w", err)
	}

	return v, nil
}

func miss(op, key string) error {
	return fmt.Errorf("%w (op: %s, key: %s)", ErrCacheMiss, op, key)
}

func wrap(op string, key string, err error) error {
	if errors.Is(err, redis.Nil) {
		return miss(op, key)
	}

	return fmt.Errorf("cache: (op: %s, key: %s): %w", op, key, err)
}
