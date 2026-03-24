package storage

import (
	"context"
	"io"
)

// FileStorage определяет интерфейс для работы с файлами.
type FileStorage interface {
	UploadFile(ctx context.Context, file io.Reader, filename string) (string, error)
	GetFileURL(ctx context.Context, fileID string) (string, error)
	DownloadFile(ctx context.Context, fileID string) ([]byte, error)
	DeleteFile(ctx context.Context, fileID string) error
}
