package cache

import (
	"context"
	"fmt"
)

func HGet[T any](ctx context.Context, m *Manager, k Key[T], id, field string) (T, error) {
	key := fullKey(m, k, id)

	raw, err := m.client.Pool.HGet(ctx, key, field).Result()
	if err != nil {
		var zero T

		return zero, wrap("HGET", key+":"+field, err)
	}

	return decode[T](raw)
}

func HSetField[T any](ctx context.Context, m *Manager, k Key[T], id, field string, v T) error {
	key := fullKey(m, k, id)

	raw, err := encode(v)
	if err != nil {
		return fmt.Errorf("cache: encode: %w", err)
	}

	if err := m.client.Pool.HSet(ctx, key, field, raw).Err(); err != nil {
		return wrap("HSET", key+":"+field, err)
	}

	return nil
}

type MapString interface {
	map[string]interface{}
}

func HSetMap[T MapString](ctx context.Context, m *Manager, k Key[T], id string, v T) error {
	key := fullKey(m, k, id)

	if err := m.client.Pool.HSet(ctx, key, v).Err(); err != nil {
		return wrap("HSET", key, err)
	}

	return nil
}
