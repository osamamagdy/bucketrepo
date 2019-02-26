package main

import (
	"context"
	"io"
	"strings"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

// CloudStorage keeps the configuration for a cloud storage
type CloudStorage struct {
	config StorageConfig
}

// NewCloudStorage creates a new cloud storage instance
func NewCloudStorage(config StorageConfig) *CloudStorage {
	return &CloudStorage{
		config: config,
	}
}

// key sanitizes the cloud storage key
func (s *CloudStorage) key(path string) string {
	return strings.TrimPrefix(path, "/")
}

// ReadFile reads a file from the cloud storage
func (s *CloudStorage) ReadFile(path string) (io.ReadCloser, error) {
	ctx, _ := context.WithTimeout(context.Background(), s.config.Timeout)
	blob, err := blob.OpenBucket(ctx, s.config.BucketURL)
	if err != nil {
		return nil, err
	}
	return blob.NewReader(ctx, s.key(path), nil)
}

// WriteFile writes a file into the cloud cloud storage
func (s *CloudStorage) WriteFile(path string, file io.ReadCloser) error {
	ctx, _ := context.WithTimeout(context.Background(), s.config.Timeout)
	blob, err := blob.OpenBucket(ctx, s.config.BucketURL)
	if err != nil {
		return err
	}
	writer, err := blob.NewWriter(ctx, s.key(path), nil)
	if err != nil {
		return nil
	}
	defer writer.Close()
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}
	return nil
}
