package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

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
	// Генерируем уникальное имя, чтобы избежать коллизий
	fileID := fmt.Sprintf("%s_%s_%d", c.cfg.UploadFolder, filename, time.Now().Unix())

	uploadParams := uploader.UploadParams{
		PublicID:     fileID,
		Folder:       c.cfg.UploadFolder,
		ResourceType: "auto",
	}

	resp, err := c.client.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return resp.PublicID, nil
}

// GetFileURL возвращает URL для скачивания файла по его ID.
func (c *CloudinaryStorage) GetFileURL(ctx context.Context, fileID string) (string, error) {
	resp, err := c.client.Upload.Explicit(ctx, uploader.ExplicitParams{
		PublicID:     fileID,
		ResourceType: "auto",
	})
	if err != nil {
		return "", fmt.Errorf("failed to get asset details via explicit: %w", err)
	}

	return resp.SecureURL, nil
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
