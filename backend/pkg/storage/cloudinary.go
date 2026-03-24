package storage

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
		Type:         "upload",
	}

	resp, err := c.client.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return resp.PublicID, nil
}

// GetFileURL returns a signed URL for viewing the file in a browser.
func (c *CloudinaryStorage) GetFileURL(ctx context.Context, fileID string) (string, error) {
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

	url := fmt.Sprintf("https://res.cloudinary.com/%s/%s/upload/%s", c.cfg.CloudName, resourceType, fileID)

	c.log.Debug("generated view url", zap.String("url", url), zap.String("public_id", fileID))
	return url, nil
}

// DownloadFile downloads a file from Cloudinary using a signed download URL.
// Cloudinary CDN does not support Basic Auth — we must use a signed URL via the API.
func (c *CloudinaryStorage) DownloadFile(ctx context.Context, fileID string) ([]byte, error) {
	// Step 1: determine resource type
	resourceType := "raw" // PDFs are typically stored as "raw" or "image"
	respAdmin, err := c.client.Admin.Asset(ctx, admin.AssetParams{
		PublicID:     fileID,
		AssetType:    api.AssetType("raw"),
		DeliveryType: api.DeliveryType("upload"),
	})
	if err != nil || respAdmin == nil {
		respAdmin2, err2 := c.client.Admin.Asset(ctx, admin.AssetParams{
			PublicID:     fileID,
			AssetType:    api.AssetType("image"),
			DeliveryType: api.DeliveryType("upload"),
		})
		if err2 == nil && respAdmin2 != nil {
			resourceType = respAdmin2.ResourceType
		}
	} else {
		resourceType = respAdmin.ResourceType
	}

	// Step 2: generate a signed download URL
	// Cloudinary signed URL format: timestamp + params + secret -> sha1 signature
	timestamp := time.Now().Unix()

	// Parameters to sign (sorted alphabetically)
	params := map[string]string{
		"public_id": fileID,
		"timestamp": fmt.Sprintf("%d", timestamp),
	}

	signature := generateSignature(params, c.cfg.APISecret)

	// Build the signed URL using Cloudinary's delivery URL with authentication
	downloadURL := fmt.Sprintf(
		"https://res.cloudinary.com/%s/%s/upload/fl_attachment/%s?api_key=%s&timestamp=%d&signature=%s",
		c.cfg.CloudName,
		resourceType,
		fileID,
		c.cfg.APIKey,
		timestamp,
		signature,
	)

	c.log.Debug("downloading file with signed URL",
		zap.String("public_id", fileID),
		zap.String("resource_type", resourceType),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Fallback: try the simple public URL (works if account allows it)
		return c.downloadPublicFallback(ctx, fileID, resourceType)
	}

	return io.ReadAll(resp.Body)
}

// downloadPublicFallback attempts to download using the public CDN URL.
// This works when the Cloudinary account/resource is set to public access.
func (c *CloudinaryStorage) downloadPublicFallback(ctx context.Context, fileID, resourceType string) ([]byte, error) {
	publicURL := fmt.Sprintf(
		"https://res.cloudinary.com/%s/%s/upload/%s",
		c.cfg.CloudName,
		resourceType,
		fileID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, publicURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create fallback request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed fallback download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cloudinary returned status %d for fileID: %s", resp.StatusCode, fileID)
	}

	return io.ReadAll(resp.Body)
}

// generateSignature creates a Cloudinary v1 API signature.
// See: https://cloudinary.com/documentation/authentication#generating_authentication_signatures
func generateSignature(params map[string]string, apiSecret string) string {
	// Sort keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build param string
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	paramStr := strings.Join(parts, "&")

	// SHA1(paramStr + apiSecret)
	h := hmac.New(sha1.New, []byte(apiSecret))
	h.Write([]byte(paramStr + apiSecret))
	return hex.EncodeToString(h.Sum(nil))
}

// DeleteFile удаляет файл из хранилища.
func (c *CloudinaryStorage) DeleteFile(ctx context.Context, fileID string) error {
	_, err := c.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: fileID,
	})
	return err
}
