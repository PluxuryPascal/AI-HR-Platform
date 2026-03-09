package grpc

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *zap.Logger
	cfg  *config.GRPCClient
	name string

	conn *grpc.ClientConn
}

// GetConn возвращает установленное соединение для передачи в конструкторы клиентов
func (c *Client) GetConn() *grpc.ClientConn {
	return c.conn
}

func (c *Client) DependsOn() []string {
	return []string{"logger"}
}

func (c *Client) HealthCheck(ctx context.Context) error {
	if c.conn == nil {
		return fmt.Errorf("grpc client is not initialized")
	}

	state := c.conn.GetState()

	if state == connectivity.Shutdown {
		return fmt.Errorf("grpc client is shut down")
	}

	if state == connectivity.TransientFailure {
		return fmt.Errorf("grpc client is in transient failure state")
	}

	c.conn.Connect()

	return nil
}

func (c *Client) Init(ctx context.Context) error {
	opts := []grpc.DialOption{}

	if c.cfg.MaxRecvMsgSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(c.cfg.MaxRecvMsgSize)))
	}

	if c.cfg.UseTLS {
		cert, err := tls.LoadX509KeyPair(c.cfg.TLS.ClientCertPath, c.cfg.TLS.ClientKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load x509 key pair: %w", err)
		}

		caCertPool := x509.NewCertPool()

		caBytes, err := os.ReadFile(c.cfg.TLS.CaCertPath)
		if err != nil {
			return fmt.Errorf("failed to read ca cert: %w", err)
		}

		if ok := caCertPool.AppendCertsFromPEM(caBytes); !ok {
			return fmt.Errorf("failed to append ca cert")
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caCertPool,
			MinVersion:   tls.VersionTLS13,
		}

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port),
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to create grpc client: %w", err)
	}

	c.conn = conn

	c.log.Info("grpc client started", zap.Int("port", c.cfg.Port))

	return nil
}

func (c *Client) Name() string {
	return fmt.Sprintf("grpc-client-%s", c.name)
}

func (c *Client) Run(ctx context.Context) error {
	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close grpc client: %w", err)
	}

	return nil
}

func NewClient(name string, log *zap.Logger, cfg *config.GRPCClient) *Client {
	return &Client{
		name: name,
		log:  log,
		cfg:  cfg,
	}
}

var _ svc.Service = (*Client)(nil)
