package cache

import (
	"context"
	"fmt"
	"time"
)

func Get[T any](ctx context.Context, m *Manager, k Key[T], id string) (T, error) {
	key := fullKey(m, k, id)

	raw, err := m.client.Pool.Get(ctx, key).Result()
	if err != nil {
		var zero T
		return zero, wrap("GET", key, err)
	}

	return decode[T](raw)
}

func SetWithTTL[T any](ctx context.Context, m *Manager, k Key[T], id string, v T, ttl time.Duration) error {
	key := fullKey(m, k, id)

	raw, err := encode(v)
	if err != nil {
		return fmt.Errorf("cache: encode: %w", err)
	}

	if err := m.client.Pool.Set(ctx, key, raw, ttl).Err(); err != nil {
		return wrap("SET", key, err)
	}

	return nil
}

func Set[T any](ctx context.Context, m *Manager, k Key[T], id string, v T) error {
	return SetWithTTL(ctx, m, k, id, v, k.ttl)
}

func Scan[T any](ctx context.Context, m *Manager, k Key[T]) (int64, error) {
	var (
		count  int64
		cursor uint64
	)

	pattern := fullKey(m, k, "*")

	for {
		keys, c, err := m.client.Pool.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return 0, wrap("SCAN", pattern, err)
		}

		count += int64(len(keys))

		cursor = c

		if cursor == 0 {
			break
		}
	}

	return count, nil
}

func Incr[T any](ctx context.Context, m *Manager, k Key[T], id string) (int64, error) {
	key := fullKey(m, k, id)

	count, err := m.client.Pool.Incr(ctx, key).Result()
	if err != nil {
		return 0, wrap("INCR", key, err)
	}

	return count, nil
}

func IncrWithTTL[T any](ctx context.Context, m *Manager, k Key[T], id string, ttl time.Duration) (int64, error) {
	key := fullKey(m, k, id)

	pipe := m.client.Pool.TxPipeline()

	incr, err := pipe.Incr(ctx, key).Result()
	if err != nil {
		return 0, wrap("INCR_WITH_TTL", key, err)
	}

	pipe.Expire(ctx, key, ttl)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, wrap("INCR_WITH_TTL", key, err)
	}

	return incr, nil
}

func Delete[T any](ctx context.Context, m *Manager, k Key[T], id string) error {
	key := fullKey(m, k, id)

	if err := m.client.Pool.Del(ctx, key).Err(); err != nil {
		return wrap("DELETE", key, err)
	}

	return nil
}
