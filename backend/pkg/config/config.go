package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server    Server               `yaml:"server"`
	Database  Database             `yaml:"database"`
	Redis     Redis                `yaml:"redis"`
	Hash      Hash                 `yaml:"hash"`
	Logger    Logger               `yaml:"logger"`
	Token     Token                `yaml:"token"`
	RateLimit map[string]RateLimit `yaml:"rate-limit"`
	Invite    Invite               `yaml:"invite"`
}

type Invite struct {
	TTL time.Duration `yaml:"ttl"`
}

type RateLimit struct {
	Requests int           `yaml:"requests"`
	Window   time.Duration `yaml:"window"`
}

type Token struct {
	Issuer     string        `yaml:"issuer"`
	ExpireAt   time.Duration `yaml:"expire-at"`
	PrivateKey PrivateKey    `yaml:"private-key"`
}

type PrivateKey struct {
	Path string `yaml:"path"`
}

type Logger struct {
	StdOut bool   `yaml:"stdout"`
	Level  string `yaml:"level"`
	File   File   `yaml:"file"`
}

type File struct {
	Path string `yaml:"path"`
}

type Hash struct {
	Time    uint32 `yaml:"time"`
	Memory  uint32 `yaml:"memory"`
	Threads uint8  `yaml:"threads"`
	KeyLen  uint32 `yaml:"key-len"`
	SaltLen uint32 `yaml:"salt-len"`
}

type Redis struct {
	Host              string        `yaml:"host"`
	ConnectionTimeout time.Duration `yaml:"connection-timeout"`
}

type Server struct {
	Port         int           `yaml:"port"`
	SecureCookie bool          `yaml:"secure-cookie"`
	ReadTimeout  time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
	IdleTimeout  time.Duration `yaml:"idle-timeout"`
}

type Database struct {
	Host              string        `yaml:"host"`
	ConnectionTimeout time.Duration `yaml:"connection-timeout"`
	MaxConns          int32         `yaml:"max-conns"`
	MinConns          int32         `yaml:"min-conns"`
	MaxConnLifetime   time.Duration `yaml:"max-conn-lifetime"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	path = filepath.Clean(path)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}
