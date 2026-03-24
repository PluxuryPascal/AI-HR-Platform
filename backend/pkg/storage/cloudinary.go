package storage

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CloudinaryStorage struct {
	log *zap.Logger
	cfg *config.Cloudinary

	client *cloudinary.Cloudinary
}

func (c *CloudinaryStorage) Init(ctx context.Context) error {
	if c.cfg.URL == "" {
		return fmt.Errorf("cloudinary url is required")
	}

	cld, err := cloudinary.NewFromURL(c.cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to create cloudinary client from URL: %w", err)
	}

	c.client = cld
	c.log.Debug("cloudinary initialized from URL", zap.String("cloud_name", c.cfg.CloudName))
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
var _ FileStorage = (*CloudinaryStorage)(nil)

func NewCloudinaryStorage(log *zap.Logger, cfg *config.Cloudinary) *CloudinaryStorage {
	return &CloudinaryStorage{
		log: log,
		cfg: cfg,
	}
}

// UploadFile загружает файл в хранилище и возвращает публичный ID (ключ).
func (c *CloudinaryStorage) UploadFile(ctx context.Context, file io.Reader, filename string) (string, error) {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	fileID := fmt.Sprintf("%s_%s", base, uuid.New().String())

	uploadParams := uploader.UploadParams{
		PublicID:     fileID,
		Folder:       c.cfg.UploadFolder,
		ResourceType: "auto",
		Type:         "upload", // Public upload for reliability
	}

	resp, err := c.client.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return resp.PublicID, nil
}

// GetFileURL возвращает URL для скачивания файла по его ID.
func (c *CloudinaryStorage) GetFileURL(ctx context.Context, fileID string) (string, error) {
	// 1. Fetch asset metadata to determine exact ResourceType
	var resourceType string

	fetchMetadata := func(resType string) bool {
		respAdmin, errAdmin := c.client.Admin.Asset(ctx, admin.AssetParams{
			PublicID:     fileID,
			AssetType:    api.AssetType(resType),
			DeliveryType: api.DeliveryType("upload"),
		})
		if errAdmin == nil && respAdmin != nil {
			resourceType = respAdmin.ResourceType
			return true
		}
		return false
	}

	found := fetchMetadata("image")
	if !found {
		found = fetchMetadata("raw")
	}

	if !found {
		return "", fmt.Errorf("cloudinary could not find metadata for fileID: %s", fileID)
	}

	// 2. Generate a direct public URL. 
	// To avoid 401 unauthorized, we use the simplest "upload" type URL.
	// If Cloudinary still returns 401, it's due to account-level strict transformation settings.
	url := fmt.Sprintf("https://res.cloudinary.com/%s/%s/upload/%s", c.cfg.CloudName, resourceType, fileID)

	c.log.Debug("generated download url", zap.String("url", url), zap.String("public_id", fileID), zap.String("resource_type", resourceType))
	return url, nil
}

// DownloadFile скачивает файл из хранилища в память (возвращает байты).
func (c *CloudinaryStorage) DownloadFile(ctx context.Context, fileID string) ([]byte, error) {
	url, err := c.GetFileURL(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// Use HTTP GET with Basic Auth to skip Any 401 issues for public/private assets
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	// Cloudinary allows Basic Auth with APIKey:APISecret
	req.SetBasicAuth(c.cfg.APIKey, c.cfg.APISecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cloudinary returned status %d for URL: %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)
}

// DeleteFile удаляет файл из хранилища.
func (c *CloudinaryStorage) DeleteFile(ctx context.Context, fileID string) error {
	_, err := c.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: fileID,
	})
	return err
}
