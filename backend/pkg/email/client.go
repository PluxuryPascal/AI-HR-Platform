package email

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"fmt"

	"github.com/wneessen/go-mail"
	"go.uber.org/zap"
)

type Client struct {
	log *zap.Logger
	cfg *config.SMTP
}

var _ svc.Service = (*Client)(nil)

func NewClient(log *zap.Logger, cfg *config.SMTP) *Client {
	return &Client{log: log, cfg: cfg}
}

func (c *Client) Name() string        { return "email" }
func (c *Client) DependsOn() []string { return []string{"logger"} }

func (c *Client) Init(_ context.Context) error {
	if c.cfg == nil {
		return fmt.Errorf("smtp config is nil")
	}

	if c.cfg.Host == "" {
		return fmt.Errorf("smtp host is empty")
	}

	c.log.Info("email client initialized",
		zap.String("host", c.cfg.Host),
		zap.Int("port", c.cfg.Port),
		zap.String("from", c.cfg.From),
		zap.Bool("use_tls", c.cfg.UseTLS),
	)

	return nil
}

func (c *Client) HealthCheck(_ context.Context) error {
	if c.cfg == nil {
		return fmt.Errorf("email client is not configured")
	}

	return nil
}

func (c *Client) Run(_ context.Context) error { return nil }

func (c *Client) Stop(_ context.Context) error { return nil }

func (c *Client) Send(ctx context.Context, to, subject, body string) error {
	m := mail.NewMsg()

	if err := m.From(c.cfg.From); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}

	if err := m.To(to); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}

	m.Subject(subject)
	m.SetBodyString(mail.TypeTextPlain, body)

	tlsPolicy := mail.NoTLS
	if c.cfg.UseTLS {
		tlsPolicy = mail.TLSMandatory
	}

	opts := []mail.Option{
		mail.WithPort(c.cfg.Port),
		mail.WithTLSPolicy(tlsPolicy),
	}

	if c.cfg.Username != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(c.cfg.Username),
			mail.WithPassword(c.cfg.Password),
		)
	}

	client, err := mail.NewClient(c.cfg.Host, opts...)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	if err := client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	c.log.Debug("email sent via go-mail", zap.String("to", to), zap.String("subject", subject))

	return nil
}
