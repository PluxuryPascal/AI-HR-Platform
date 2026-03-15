package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server         Server               `yaml:"server"`
	Database       Database             `yaml:"database"`
	Redis          Redis                `yaml:"redis"`
	Hash           Hash                 `yaml:"hash"`
	Logger         Logger               `yaml:"logger"`
	Token          Token                `yaml:"token"`
	RateLimit      map[string]RateLimit `yaml:"rate-limit"`
	Invite         Invite               `yaml:"invite"`
	InviteRecovery InviteRecovery       `yaml:"invite-recovery"`
	RabbitMQ       RabbitMQ             `yaml:"rabbitmq"`
	Cloudinary     Cloudinary           `yaml:"cloudinary"`
	GRPC           GRPC                 `yaml:"grpc"`
	Temporal       Temporal             `yaml:"temporal"`
}

type Temporal struct {
	HostPort        string        `yaml:"host-port"`
	ActivityTimeout time.Duration `yaml:"activity-timeout"`
	QueueName       string        `yaml:"queue-name"`
	Namespace       string        `yaml:"namespace"`
	WorkerCount     int           `yaml:"worker-count"`
}

type Invite struct {
	TTL time.Duration `yaml:"ttl"`
}

type InviteRecovery struct {
	Cron           string        `yaml:"cron"`
	StuckThreshold time.Duration `yaml:"stuck-threshold"`
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
	Ports        map[string]int `yaml:"ports"`
	SecureCookie bool           `yaml:"secure-cookie"`
	ReadTimeout  time.Duration  `yaml:"read-timeout"`
	WriteTimeout time.Duration  `yaml:"write-timeout"`
	IdleTimeout  time.Duration  `yaml:"idle-timeout"`
}

type Database struct {
	Host              string        `yaml:"host"`
	ConnectionTimeout time.Duration `yaml:"connection-timeout"`
	MaxConns          int32         `yaml:"max-conns"`
	MinConns          int32         `yaml:"min-conns"`
	MaxConnLifetime   time.Duration `yaml:"max-conn-lifetime"`
}

type ExchangeConfig struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Durable bool   `yaml:"durable"`
}

type QueueConfig struct {
	Name          string `yaml:"name"`
	Exchange      string `yaml:"exchange"`
	RoutingKey    string `yaml:"routing_key"`
	DLX           string `yaml:"dlx"`
	DLXRoutingKey string `yaml:"dlx_routing_key"`
	MaxRetries    int    `yaml:"max_retries"`
	MessageTTL    int    `yaml:"message_ttl"`
	PrefetchCount int    `yaml:"prefetch_count"`
	Concurrency   int    `yaml:"concurrency"`
}

type RabbitMQ struct {
	URL            string           `yaml:"url"`
	Exchanges      []ExchangeConfig `yaml:"exchanges"`
	Queues         []QueueConfig    `yaml:"queues"`
	ReconnectDelay time.Duration    `yaml:"reconnect_delay"`
}

type Cloudinary struct {
	URL          string `yaml:"url"`
	CloudName    string `yaml:"cloud_name"`
	APIKey       string `yaml:"api_key"`
	APISecret    string `yaml:"api_secret"`
	UploadFolder string `yaml:"upload_folder"`
}

type GRPC struct {
	Servers map[string]GRPCServer `yaml:"servers"`
	Clients map[string]GRPCClient `yaml:"clients"`
}

type GRPCServer struct {
	Port           int  `yaml:"port"`
	MaxRecvMsgSize int  `yaml:"max_recv_msg_size"`
	UseTLS         bool `yaml:"use_tls"`
	TLS            TLS  `yaml:"tls"`
}

type GRPCClient struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	MaxRecvMsgSize int    `yaml:"max_recv_msg_size"`
	UseTLS         bool   `yaml:"use_tls"`
	TLS            TLS    `yaml:"tls"`
}

type TLS struct {
	CaCertPath     string `yaml:"ca_cert_path"`
	ServerCertPath string `yaml:"server_cert_path"`
	ServerKeyPath  string `yaml:"server_key_path"`
	ClientCertPath string `yaml:"client_cert_path"`
	ClientKeyPath  string `yaml:"client_key_path"`
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
