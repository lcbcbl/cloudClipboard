package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"cloud-clipboard/app/api"
	"cloud-clipboard/app/config"
	"cloud-clipboard/internal/clipboard"
	"cloud-clipboard/internal/file"
	"cloud-clipboard/internal/logger"
)

func main() {
	// 加载配置
	cfg := config.GetDefaultConfig()

	// 初始化日志
	if err := logger.InitLogger(logger.GetDefaultConfig()); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}

	// 初始化服务
	cache := clipboard.NewLRUCache(cfg.Clipboard.MaxMemory, cfg.Clipboard.MaxItems)

	fileService, err := file.NewFileService(cfg.File.UploadDir, cfg.File.MetadataFile)
	if err != nil {
		logger.Fatalf("Failed to initialize file service: %v", err)
	}

	// 初始化控制器
	clipboardController := api.NewClipboardController(cache, &cfg.Clipboard)
	fileController := api.NewFileController(fileService, &cfg.File)

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 静态文件服务
	r.Static("/uploads", cfg.File.UploadDir)

	// API路由
	api := r.Group("/api")
	{
		// 字符串剪切板路由
		clipboard := api.Group("/clipboard")
		{
			clipboard.POST("/text", clipboardController.UploadText)
			clipboard.GET("/text", clipboardController.GetAllText)
			clipboard.DELETE("/text", clipboardController.ClearAllText)
			clipboard.GET("/text/:id", clipboardController.GetTextById)
			clipboard.DELETE("/text/:id", clipboardController.DeleteTextById)
		}

		// 文件路由
		files := api.Group("/files")
		{
			files.POST("", fileController.UploadFile)
			files.GET("", fileController.GetAllFiles)
			files.GET("/:id", fileController.GetFileInfo)
			files.GET("/:id/download", fileController.DownloadFile)
			files.GET("/:id/thumbnail", fileController.GetFileThumbnail)
			files.DELETE("/:id", fileController.DeleteFile)
		}
	}

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 设置定期清理任务
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.File.CleanupInterval) * time.Millisecond)
		defer ticker.Stop()

		for {
			<-ticker.C
			logger.Info("Running file cleanup task...")
			deletedCount, err := fileService.CleanupExpiredFiles(cfg.File.MaxAge)
			if err != nil {
				logger.Errorf("Failed to cleanup expired files: %v", err)
				continue
			}
			logger.Infof("Cleanup completed. Deleted %d expired files.", deletedCount)
		}
	}()

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Server is running on http://%s", addr)
	logger.Info("API endpoints:")
	logger.Info("  POST   /api/clipboard/text      - Upload text")
	logger.Info("  GET    /api/clipboard/text      - Get all text items")
	logger.Info("  GET    /api/clipboard/text/:id  - Get specific text item")
	logger.Info("  DELETE /api/clipboard/text/:id  - Delete text item")
	logger.Info("  DELETE /api/clipboard/text      - Clear all text items")
	logger.Info("  POST   /api/files               - Upload file")
	logger.Info("  GET    /api/files               - Get all files")
	logger.Info("  GET    /api/files/:id           - Get file info")
	logger.Info("  GET    /api/files/:id/download  - Download file")
	logger.Info("  DELETE /api/files/:id           - Delete file")

	if err := r.Run(addr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
