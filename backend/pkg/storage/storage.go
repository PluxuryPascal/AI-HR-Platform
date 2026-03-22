package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// FileStorage определяет интерфейс для работы с файлами.
type FileStorage interface {
	UploadFile(ctx context.Context, file io.Reader, filename string) (string, error)
	GetFileURL(ctx context.Context, fileID string) (string, error)
	DownloadFile(ctx context.Context, fileID string) ([]byte, error)
	DeleteFile(ctx context.Context, fileID string) error
}

var _ FileStorage = (*CloudinaryStorage)(nil)

// UploadFile загружает файл в хранилище и возвращает публичный ID (ключ).
func (c *CloudinaryStorage) UploadFile(ctx context.Context, file io.Reader, filename string) (string, error) {
	// Strip extension from filename for cleaner PublicID
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	// Generate a unique public ID without the folder prefix (Folder param handles it)
	fileID := fmt.Sprintf("%s_%d", base, time.Now().Unix())

	uploadParams := uploader.UploadParams{
		PublicID:     fileID,
		Folder:       c.cfg.UploadFolder,
		ResourceType: "auto",
		Type:         "authenticated", // Use authenticated for better security
	}

	resp, err := c.client.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return resp.PublicID, nil
}

// GetFileURL возвращает URL для скачивания файла по его ID.
func (c *CloudinaryStorage) GetFileURL(ctx context.Context, fileID string) (string, error) {
	// 1. Fetch asset metadata to determine exact ResourceType and Format
	var resourceType, format string

	fetchMetadata := func(resType string) bool {
		respAdmin, errAdmin := c.client.Admin.Asset(ctx, admin.AssetParams{
			PublicID:     fileID,
			AssetType:    api.AssetType(resType),
			DeliveryType: api.DeliveryType("authenticated"),
		})
		if errAdmin == nil && respAdmin != nil && respAdmin.Format != "" {
			resourceType = respAdmin.ResourceType
			format = respAdmin.Format
			return true
		}
		return false
	}

	found := fetchMetadata("image")
	if !found {
		found = fetchMetadata("raw")
	}

	// 2. Fallback to Search API if Admin API fails
	if !found {
		searchQuery := search.Query{
			Expression: fmt.Sprintf("public_id:%s", fileID),
		}
		respSearch, errSearch := c.client.Admin.Search(ctx, searchQuery)
		if errSearch == nil && respSearch != nil && len(respSearch.Assets) > 0 {
			resourceType = respSearch.Assets[0].ResourceType
			format = respSearch.Assets[0].Format
			found = true
		}
	}

	if !found {
		return "", fmt.Errorf("cloudinary could not find metadata for fileID: %s", fileID)
	}

	// 3. Generate a Private Download URL allowing backend download of authenticated assets
	downloadURL, err := c.client.Upload.PrivateDownloadURL(uploader.PrivateDownloadURLParams{
		PublicID:     fileID,
		Format:       format,
		DeliveryType: "authenticated",
		ResourceType: api.AssetType(resourceType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate private download url: %w", err)
	}

	return downloadURL, nil
}

// DownloadFile скачивает файл из хранилища в память (возвращает байты).
func (c *CloudinaryStorage) DownloadFile(ctx context.Context, fileID string) ([]byte, error) {
	url, err := c.GetFileURL(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file url for download: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from cloudinary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cloudinary returned non-200 status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

// DeleteFile удаляет файл из хранилища.
func (c *CloudinaryStorage) DeleteFile(ctx context.Context, fileID string) error {
	_, err := c.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: fileID,
	})
	if err != nil {
		return fmt.Errorf("cloudinary delete failed: %w", err)
	}

	return nil
}
