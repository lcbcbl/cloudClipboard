package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// FileService 文件服务
type FileService struct {
	metadataPath string
	uploadDir    string
	mu           sync.RWMutex
}

// FileMetadata 文件元数据
type FileMetadata struct {
	ID             string `json:"id"`
	Filename       string `json:"filename"`
	Size           int64  `json:"size"`
	Mimetype       string `json:"mimetype"`
	FilePath       string `json:"filePath"`
	UploadTime     int64  `json:"uploadTime"`
	LastAccessTime int64  `json:"lastAccessTime"`
	DownloadCount  int    `json:"downloadCount"`
	MaxDownloads   int    `json:"maxDownloads"`
}

// NewFileService 创建新的文件服务
func NewFileService(uploadDir, metadataFile string) (*FileService, error) {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 确保元数据文件目录存在
	metadataDir := filepath.Dir(metadataFile)
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// 确保元数据文件存在
	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		if err := os.WriteFile(metadataFile, []byte("[]"), 0644); err != nil {
			return nil, fmt.Errorf("failed to create metadata file: %w", err)
		}
	}

	return &FileService{
		metadataPath: metadataFile,
		uploadDir:    uploadDir,
	}, nil
}

// ReadMetadata 读取元数据
func (s *FileService) ReadMetadata() ([]*FileMetadata, error) {
	data, err := os.ReadFile(s.metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata []*FileMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// 确保返回空切片而不是nil
	if metadata == nil {
		return []*FileMetadata{}, nil
	}

	return metadata, nil
}

// WriteMetadata 写入元数据
func (s *FileService) WriteMetadata(metadata []*FileMetadata) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(s.metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// AddFileMetadata 添加文件元数据
func (s *FileService) AddFileMetadata(fileInfo *FileInfo) (*FileMetadata, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metadata, err := s.ReadMetadata()
	if err != nil {
		return nil, err
	}

	newFile := &FileMetadata{
		ID:             uuid.New().String(),
		Filename:       fileInfo.OriginalName,
		Size:           fileInfo.Size,
		Mimetype:       fileInfo.Mimetype,
		FilePath:       fileInfo.Path,
		UploadTime:     time.Now().UnixMilli(),
		LastAccessTime: time.Now().UnixMilli(),
		DownloadCount:  0,
		MaxDownloads:   fileInfo.MaxDownloads,
	}

	metadata = append(metadata, newFile)
	if err := s.WriteMetadata(metadata); err != nil {
		return nil, err
	}

	return newFile, nil
}

// GetFileMetadata 获取文件元数据
func (s *FileService) GetFileMetadata(id string) (*FileMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metadata, err := s.ReadMetadata()
	if err != nil {
		return nil, err
	}

	for _, file := range metadata {
		if file.ID == id {
			return file, nil
		}
	}

	return nil, ErrFileNotFound
}

// GetAllFileMetadata 获取所有文件元数据
func (s *FileService) GetAllFileMetadata() ([]*FileMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ReadMetadata()
}

// UpdateFileMetadata 更新文件元数据
func (s *FileService) UpdateFileMetadata(id string, updates map[string]interface{}) (*FileMetadata, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metadata, err := s.ReadMetadata()
	if err != nil {
		return nil, err
	}

	var file *FileMetadata
	var index int

	for i, f := range metadata {
		if f.ID == id {
			file = f
			index = i
			break
		}
	}

	if file == nil {
		return nil, ErrFileNotFound
	}

	// 更新字段
	for key, value := range updates {
		switch key {
		case "downloadCount":
			if v, ok := value.(int); ok {
				file.DownloadCount = v
			}
		case "lastAccessTime":
			if v, ok := value.(int64); ok {
				file.LastAccessTime = v
			}
		}
	}

	// 更新最后访问时间
	file.LastAccessTime = time.Now().UnixMilli()

	metadata[index] = file
	if err := s.WriteMetadata(metadata); err != nil {
		return nil, err
	}

	return file, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	metadata, err := s.ReadMetadata()
	if err != nil {
		return err
	}

	var file *FileMetadata
	var index int

	for i, f := range metadata {
		if f.ID == id {
			file = f
			index = i
			break
		}
	}

	if file == nil {
		return ErrFileNotFound
	}

	// 删除实际文件
	if err := os.Remove(file.FilePath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// 删除元数据
	metadata = append(metadata[:index], metadata[index+1:]...)
	if err := s.WriteMetadata(metadata); err != nil {
		return err
	}

	return nil
}

// CleanupExpiredFiles 清理过期文件
func (s *FileService) CleanupExpiredFiles(maxAge int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metadata, err := s.ReadMetadata()
	if err != nil {
		return 0, err
	}

	now := time.Now().UnixMilli()
	var deletedCount int
	var remainingMetadata []*FileMetadata

	for _, file := range metadata {
		if now-file.UploadTime > maxAge {
			// 删除实际文件
			if err := os.Remove(file.FilePath); err != nil && !errors.Is(err, fs.ErrNotExist) {
				// 记录错误但继续执行
				fmt.Printf("Failed to delete expired file %s: %v\n", file.ID, err)
			}
			deletedCount++
		} else {
			remainingMetadata = append(remainingMetadata, file)
		}
	}

	if err := s.WriteMetadata(remainingMetadata); err != nil {
		return 0, err
	}

	return deletedCount, nil
}

// CheckTotalStorage 检查总存储大小
func (s *FileService) CheckTotalStorage() (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metadata, err := s.ReadMetadata()
	if err != nil {
		return 0, err
	}

	var totalSize int64
	for _, file := range metadata {
		totalSize += file.Size
	}

	return totalSize, nil
}

// FileInfo 文件信息
type FileInfo struct {
	OriginalName string
	Size         int64
	Mimetype     string
	Path         string
	MaxDownloads int
}

// 错误定义
var (
	ErrFileNotFound        = errors.New("file not found")
	ErrMaxDownloadsReached = errors.New("maximum download limit reached")
)
