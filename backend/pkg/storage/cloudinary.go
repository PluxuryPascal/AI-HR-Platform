package storage

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"go.uber.org/zap"
)

type CloudinaryStorage struct {
	log *zap.Logger
	cfg *config.Cloudinary

	client *cloudinary.Cloudinary
}

func (c *CloudinaryStorage) Init(ctx context.Context) error {
	cloudinary, err := cloudinary.NewFromURL(c.cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to create cloudinary client: %w", err)
	}

	c.client = cloudinary

	c.log.Debug("cloudinary initialized")

	return nil
}

func (c *CloudinaryStorage) DependsOn() []string {
	return []string{"logger"}
}

func (c *CloudinaryStorage) HealthCheck(ctx context.Context) error {
	_, err := c.client.Admin.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping cloudinary: %w", err)
	}

	c.log.Debug("cloudinary health check passed")

	return nil
}

func (c *CloudinaryStorage) Name() string {
	return "cloudinary"
}

func (c *CloudinaryStorage) Run(ctx context.Context) error {
	return nil
}
func (c *CloudinaryStorage) Stop(ctx context.Context) error {
	return nil
}

var _ svc.Service = (*CloudinaryStorage)(nil)

func NewCloudinaryStorage(log *zap.Logger, cfg *config.Cloudinary) *CloudinaryStorage {
	return &CloudinaryStorage{
		log: log,
		cfg: cfg,
	}
}
