package config

// Config 应用配置
type Config struct {
	Server    ServerConfig    `json:"server"`
	Clipboard ClipboardConfig `json:"clipboard"`
	File      FileConfig      `json:"file"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// ClipboardConfig 字符串剪切板配置
type ClipboardConfig struct {
	MaxMemory   int64 `json:"maxMemory"`
	MaxItems    int   `json:"maxItems"`
	MaxItemSize int64 `json:"maxItemSize"`
}

// FileConfig 文件配置
type FileConfig struct {
	UploadDir       string `json:"uploadDir"`
	MetadataFile    string `json:"metadataFile"`
	MaxFileSize     int64  `json:"maxFileSize"`
	MaxStorage      int64  `json:"maxStorage"`
	MaxDownloads    int    `json:"maxDownloads"`
	SpeedLimit      int64  `json:"speedLimit"`
	CleanupInterval int64  `json:"cleanupInterval"`
	MaxAge          int64  `json:"maxAge"`
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "3000",
			Host: "localhost",
		},
		Clipboard: ClipboardConfig{
			MaxMemory:   1 * 1024 * 1024, // 1MB
			MaxItems:    512,
			MaxItemSize: 1 * 1024, // 1KB
		},
		File: FileConfig{
			UploadDir:       "./uploads",
			MetadataFile:    "./data/files.json",
			MaxFileSize:     16 * 1024 * 1024,  // 16MB
			MaxStorage:      512 * 1024 * 1024, // 512GB
			MaxDownloads:    10,
			SpeedLimit:      1 * 1024 * 1024,         // 1MB/s
			CleanupInterval: 24 * 60 * 60 * 1000,     // 24小时
			MaxAge:          7 * 24 * 60 * 60 * 1000, // 7天
		},
	}
}
