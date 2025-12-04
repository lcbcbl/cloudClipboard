package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"cloud-clipboard/app/config"
	"cloud-clipboard/internal/errors"
	fileservice "cloud-clipboard/internal/file"
	"cloud-clipboard/internal/logger"
)

// FileController 文件控制器
type FileController struct {
	fileService *fileservice.FileService
	config      *config.FileConfig
}

// NewFileController 创建新的文件控制器
func NewFileController(fileService *fileservice.FileService, config *config.FileConfig) *FileController {
	return &FileController{
		fileService: fileService,
		config:      config,
	}
}

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传文件到服务器
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的文件"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files [post]
func (c *FileController) UploadFile(ctx *gin.Context) {
	// 获取上传的文件，不使用http.MaxBytesReader，因为它会关闭连接
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		logger.Errorf("Failed to get file from request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeFileSizeExceeded,
			"message": fmt.Sprintf("文件大小超过限制（最大%vMB）", c.config.MaxFileSize/(1024*1024)),
		})
		return
	}
	defer file.Close()

	// 检查文件大小
	if header.Size > c.config.MaxFileSize {
		logger.Warnf("File size exceeds maximum limit: %d, max allowed: %d", header.Size, c.config.MaxFileSize)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeFileSizeExceeded,
			"message": fmt.Sprintf("文件大小超过限制（最大%vMB）", c.config.MaxFileSize/(1024*1024)),
		})
		return
	}

	// 检查总存储限制
	totalStorage, err := c.fileService.CheckTotalStorage()
	if err != nil {
		logger.Errorf("Failed to check total storage: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeCheckStorageFailed,
			"message": "检查总存储大小失败",
		})
		return
	}

	if totalStorage+header.Size > c.config.MaxStorage {
		logger.Warnf("Total storage limit exceeded. Current: %d, Max: %d, New file: %d", totalStorage, c.config.MaxStorage, header.Size)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeTotalStorageExceeded,
			"message": "总存储容量超过限制",
		})
		return
	}

	// 创建唯一文件名
	uniqueFilename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), header.Filename)
	filePath := filepath.Join(c.config.UploadDir, uniqueFilename)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("Failed to create file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeCreateFileFailed,
			"message": "创建文件失败",
		})
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		logger.Errorf("Failed to save file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeSaveFileFailed,
			"message": "保存文件内容失败",
		})
		return
	}

	// 添加文件元数据
	fileInfo := &fileservice.FileInfo{
		OriginalName: header.Filename,
		Size:         header.Size,
		Mimetype:     header.Header.Get("Content-Type"),
		Path:         filePath,
		MaxDownloads: c.config.MaxDownloads,
	}

	metadata, err := c.fileService.AddFileMetadata(fileInfo)
	if err != nil {
		logger.Errorf("Failed to add file metadata: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeAddMetadataFailed,
			"message": "添加文件元数据失败",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "File uploaded successfully",
		"file": map[string]interface{}{
			"id":         metadata.ID,
			"filename":   metadata.Filename,
			"size":       metadata.Size,
			"mimetype":   metadata.Mimetype,
			"uploadTime": metadata.UploadTime,
		},
	})
}

// GetAllFiles 获取所有文件
// @Summary 获取所有文件
// @Description 获取所有文件列表
// @Tags files
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files [get]
func (c *FileController) GetAllFiles(ctx *gin.Context) {
	files, err := c.fileService.GetAllFileMetadata()
	if err != nil {
		logger.Errorf("Failed to get all files: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeGetFilesFailed,
			"message": "获取文件列表失败",
		})
		return
	}

	// 转换为前端需要的格式，确保result始终是切片而非nil
	result := make([]map[string]interface{}, 0)
	for _, file := range files {
		result = append(result, map[string]interface{}{
			"id":             file.ID,
			"filename":       file.Filename,
			"size":           file.Size,
			"mimetype":       file.Mimetype,
			"uploadTime":     file.UploadTime,
			"lastAccessTime": file.LastAccessTime,
			"downloadCount":  file.DownloadCount,
			"maxDownloads":   file.MaxDownloads,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"files": result,
	})
}

// GetFileInfo 获取文件信息
// @Summary 获取文件信息
// @Description 根据ID获取文件信息
// @Tags files
// @Produce json
// @Param id path string true "文件ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files/{id} [get]
func (c *FileController) GetFileInfo(ctx *gin.Context) {
	id := ctx.Param("id")

	file, err := c.fileService.GetFileMetadata(id)
	if err != nil {
		if err == fileservice.ErrFileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    errors.ErrCodeFileNotFound,
				"message": "文件不存在",
			})
			return
		}
		logger.Errorf("Failed to get file info: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeGetFileInfoFailed,
			"message": "获取文件信息失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":             file.ID,
		"filename":       file.Filename,
		"size":           file.Size,
		"mimetype":       file.Mimetype,
		"uploadTime":     file.UploadTime,
		"lastAccessTime": file.LastAccessTime,
		"downloadCount":  file.DownloadCount,
		"maxDownloads":   file.MaxDownloads,
	})
}

// DownloadFile 下载文件
// @Summary 下载文件
// @Description 根据ID下载文件（带速度限制）
// @Tags files
// @Produce octet-stream
// @Param id path string true "文件ID"
// @Success 200 {file} file
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files/{id}/download [get]
func (c *FileController) DownloadFile(ctx *gin.Context) {
	id := ctx.Param("id")

	// 获取文件元数据
	file, err := c.fileService.GetFileMetadata(id)
	if err != nil {
		if err == fileservice.ErrFileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    40401,
				"message": "文件不存在",
			})
			return
		}
		logger.Errorf("Failed to get file metadata for download: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeGetFileMetaFailed,
			"message": "获取文件信息失败",
		})
		return
	}

	// 检查下载次数
	if file.DownloadCount >= file.MaxDownloads {
		logger.Warnf("File download limit reached: %s, current: %d, max: %d", id, file.DownloadCount, file.MaxDownloads)
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    errors.ErrCodeDownloadLimitReached,
			"message": "文件下载次数已达上限",
		})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		// 文件不存在，清理元数据
		logger.Warnf("File not found on disk, cleaning metadata: %s", id)
		c.fileService.DeleteFile(id)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    errors.ErrCodeFileDeleted,
			"message": "文件已被删除",
		})
		return
	}

	// 更新下载次数
	_, err = c.fileService.UpdateFileMetadata(id, map[string]interface{}{
		"downloadCount": file.DownloadCount + 1,
	})
	if err != nil {
		logger.Errorf("Failed to update download count: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeUpdateDownloadCountFailed,
			"message": "更新下载次数失败",
		})
		return
	}

	// 设置响应头
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	ctx.Header("Content-Type", file.Mimetype)
	ctx.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	// 打开文件
	src, err := os.Open(file.FilePath)
	if err != nil {
		logger.Errorf("Failed to open file for download: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeOpenFileFailed,
			"message": "打开文件失败",
		})
		return
	}
	defer src.Close()

	// 实现速度限制的文件传输
	c.speedLimitedCopy(ctx.Writer, src, file.Size, c.config.SpeedLimit)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 根据ID删除文件
// @Tags files
// @Produce json
// @Param id path string true "文件ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files/{id} [delete]
func (c *FileController) DeleteFile(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.fileService.DeleteFile(id); err != nil {
		if err == fileservice.ErrFileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    40401,
				"message": "文件不存在",
			})
			return
		}
		logger.Errorf("Failed to delete file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeDeleteFileFailed,
			"message": "删除文件失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "文件删除成功",
	})
}

// GetFileThumbnail 获取文件缩略图
// @Summary 获取文件缩略图
// @Description 根据ID获取文件缩略图，仅支持图片文件
// @Tags files
// @Produce octet-stream
// @Param id path string true "文件ID"
// @Success 200 {file} file
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files/{id}/thumbnail [get]
func (c *FileController) GetFileThumbnail(ctx *gin.Context) {
	id := ctx.Param("id")

	// 获取文件元数据
	file, err := c.fileService.GetFileMetadata(id)
	if err != nil {
		if err == fileservice.ErrFileNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    40401,
				"message": "文件不存在",
			})
			return
		}
		logger.Errorf("Failed to get file metadata for thumbnail: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrCodeGetFileMetaFailed,
			"message": "获取文件信息失败",
		})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		// 文件不存在，清理元数据
		logger.Warnf("File not found on disk, cleaning metadata: %s", id)
		c.fileService.DeleteFile(id)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    errors.ErrCodeFileDeleted,
			"message": "文件已被删除",
		})
		return
	}

	// 检查文件是否为图片
	isImage := false
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	fileExt := filepath.Ext(file.FilePath)
	if imageExtensions[fileExt] {
		isImage = true
	}

	if !isImage {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidFileFormat,
			"message": "该文件类型不支持缩略图",
		})
		return
	}

	// 设置响应头
	ctx.Header("Content-Type", file.Mimetype)
	ctx.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	// 返回文件内容
	ctx.File(file.FilePath)
}

// speedLimitedCopy 带速度限制的文件复制
func (c *FileController) speedLimitedCopy(dst io.Writer, src io.Reader, size, speedLimit int64) {
	buffer := make([]byte, 64*1024) // 64KB缓冲区
	var totalWritten int64
	startTime := time.Now()

	for {
		// 计算剩余时间和剩余数据
		elapsed := time.Since(startTime).Milliseconds()
		expectedTime := (totalWritten * 1000) / speedLimit
		if elapsed < expectedTime {
			// 需要延迟
			time.Sleep(time.Duration(expectedTime-elapsed) * time.Millisecond)
		}

		// 计算本次可以读取的数据量
		remaining := speedLimit - (totalWritten % speedLimit)
		if remaining > int64(len(buffer)) {
			remaining = int64(len(buffer))
		}

		// 读取数据
		n, err := src.Read(buffer[:remaining])
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading file: %v\n", err)
			}
			break
		}

		// 写入数据
		n, err = dst.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			break
		}

		totalWritten += int64(n)
		if totalWritten >= size {
			break
		}
	}
}
