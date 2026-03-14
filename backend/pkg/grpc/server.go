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
	log  *zap.Logger
	cfg  *config.GRPCServer
	name string

	lis    net.Listener
	server *grpc.Server

	onInit []func(*Server)
}

// GetServer возвращает инстанс gRPC сервера для регистрации на нём хендлеров
func (s *Server) GetServer() *grpc.Server {
	return s.server
}

// OnInit регистрирует функцию, которая будет вызвана в конце Init(),
// когда gRPC сервер уже создан и готов к регистрации хэндлеров.
func (s *Server) OnInit(fn func(*Server)) {
	s.onInit = append(s.onInit, fn)
}

func (s *Server) DependsOn() []string {
	return []string{"logger"}
}

func (s *Server) HealthCheck(ctx context.Context) error {
	if s.server == nil {
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
	s.server = server

	for _, fn := range s.onInit {
		fn(s)
	}

	s.log.Info("grpc server started", zap.Int("port", s.cfg.Port))

	return nil
}

func (s *Server) Name() string {
	return fmt.Sprintf("grpc-server-%s", s.name)
}

func (s *Server) Run(ctx context.Context) error {
	if err := s.server.Serve(s.lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		s.server.GracefulStop()
	}

	return nil
}

func NewServer(name string, log *zap.Logger, cfg *config.GRPCServer) *Server {
	return &Server{
		name: name,
		log:  log,
		cfg:  cfg,
	}
}

var _ svc.Service = (*Server)(nil)
