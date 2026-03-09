package grpc

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	log *zap.Logger
	cfg *config.GRPCServer

	lis    net.Listener
	Server *grpc.Server
}

// GetServer возвращает инстанс gRPC сервера для регистрации на нём хендлеров
func (s *Server) GetServer() *grpc.Server {
	return s.Server
}

func (s *Server) DependsOn() []string {
	return []string{"logger"}
}

func (s *Server) HealthCheck(ctx context.Context) error {
	if s.Server == nil {
		return fmt.Errorf("grpc server is not initialized")
	}

	if s.lis == nil {
		return fmt.Errorf("grpc server listener is not bound")
	}

	return nil
}

func (s *Server) Init(ctx context.Context) error {
	opts := []grpc.ServerOption{}

	if s.cfg.MaxRecvMsgSize > 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(s.cfg.MaxRecvMsgSize))
	}

	if s.cfg.UseTLS {
		cert, err := tls.LoadX509KeyPair(s.cfg.TLS.ServerCertPath, s.cfg.TLS.ServerKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load x509 key pair: %w", err)
		}

		caCertPool := x509.NewCertPool()

		caBytes, err := os.ReadFile(s.cfg.TLS.CaCertPath)
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

		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	server := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on %d port: %w", s.cfg.Port, err)
	}

	s.lis = lis
	s.Server = server

	s.log.Info("grpc server started", zap.Int("port", s.cfg.Port))

	return nil
}

func (s *Server) Name() string {
	return "grpc-server"
}

func (s *Server) Run(ctx context.Context) error {
	if err := s.Server.Serve(s.lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.Server != nil {
		s.Server.GracefulStop()
	}

	return nil
}

func NewServer(log *zap.Logger, cfg *config.GRPCServer) *Server {
	return &Server{
		log: log,
		cfg: cfg,
	}
}

var _ svc.Service = (*Server)(nil)
