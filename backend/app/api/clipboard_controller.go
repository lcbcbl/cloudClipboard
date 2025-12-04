package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloud-clipboard/app/config"
	"cloud-clipboard/internal/clipboard"
)

// ClipboardController 字符串剪切板控制器
type ClipboardController struct {
	cache  *clipboard.LRUCache
	config *config.ClipboardConfig
}

// NewClipboardController 创建新的字符串剪切板控制器
func NewClipboardController(cache *clipboard.LRUCache, config *config.ClipboardConfig) *ClipboardController {
	return &ClipboardController{
		cache:  cache,
		config: config,
	}
}

// UploadTextRequest 上传字符串请求
type UploadTextRequest struct {
	Text string `json:"text" binding:"required"`
}

// UploadText 上传字符串
// @Summary 上传字符串
// @Description 上传字符串到剪切板
// @Tags clipboard
// @Accept json
// @Produce json
// @Param text body UploadTextRequest true "要上传的字符串"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/clipboard/text [post]
func (c *ClipboardController) UploadText(ctx *gin.Context) {
	var req UploadTextRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid text data",
		})
		return
	}

	// 检查单条字符串大小限制
	if int64(len([]byte(req.Text))) > c.config.MaxItemSize {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Text size exceeds maximum limit",
		})
		return
	}

	id := uuid.New().String()
	if err := c.cache.Put(id, req.Text); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"text":    req.Text,
		"size":    len([]byte(req.Text)),
		"message": "Text uploaded successfully",
	})
}

// GetAllText 获取所有字符串
// @Summary 获取所有字符串
// @Description 获取所有字符串（按最近访问排序）
// @Tags clipboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/clipboard/text [get]
func (c *ClipboardController) GetAllText(ctx *gin.Context) {
	items := c.cache.GetAll()

	ctx.JSON(http.StatusOK, gin.H{
		"items":      items,
		"totalSize":  c.cache.GetSize(),
		"totalItems": c.cache.GetCount(),
	})
}

// GetTextById 获取指定字符串
// @Summary 获取指定字符串
// @Description 根据ID获取指定字符串
// @Tags clipboard
// @Produce json
// @Param id path string true "字符串ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/clipboard/text/{id} [get]
func (c *ClipboardController) GetTextById(ctx *gin.Context) {
	id := ctx.Param("id")

	text, ok := c.cache.Get(id)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Text not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":   id,
		"text": text,
	})
}

// DeleteTextById 删除指定字符串
// @Summary 删除指定字符串
// @Description 根据ID删除指定字符串
// @Tags clipboard
// @Produce json
// @Param id path string true "字符串ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/clipboard/text/{id} [delete]
func (c *ClipboardController) DeleteTextById(ctx *gin.Context) {
	id := ctx.Param("id")

	if !c.cache.Delete(id) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Text not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Text deleted successfully",
	})
}

// ClearAllText 清空所有字符串
// @Summary 清空所有字符串
// @Description 清空剪切板中的所有字符串
// @Tags clipboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/clipboard/text [delete]
func (c *ClipboardController) ClearAllText(ctx *gin.Context) {
	c.cache.Clear()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "All text items cleared successfully",
	})
}
