package cache

import (
	"backend/internal/db"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache: key not found")

type Manager struct {
	client     *db.RedisClient
	defaultTTL time.Duration
	prefix     string
}

type Option func(*Manager)

func WithDefaultTTL(ttl time.Duration) Option {
	return func(m *Manager) {
		m.defaultTTL = ttl
	}
}

func WithPrefix(prefix string) Option {
	return func(m *Manager) {
		m.prefix = prefix
	}
}

func NewManager(client *db.RedisClient, opts ...Option) *Manager {
	m := &Manager{
		client:     client,
		defaultTTL: 0,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}
