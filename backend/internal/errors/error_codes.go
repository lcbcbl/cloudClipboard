package errors

// 错误码常量定义

// 400 Bad Request
const (
	// ErrCodeFileSizeExceeded 文件大小超过限制
	ErrCodeFileSizeExceeded = 40001
	// ErrCodeTotalStorageExceeded 总存储容量超过限制
	ErrCodeTotalStorageExceeded = 40002
	// ErrCodeInvalidFileFormat 文件格式无效
	ErrCodeInvalidFileFormat = 40003
)

// 403 Forbidden
const (
	// ErrCodeDownloadLimitReached 文件下载次数已达上限
	ErrCodeDownloadLimitReached = 40301
)

// 404 Not Found
const (
	// ErrCodeFileNotFound 文件不存在
	ErrCodeFileNotFound = 40401
	// ErrCodeFileDeleted 文件已被删除
	ErrCodeFileDeleted = 40402
)

// 500 Internal Server Error
const (
	// ErrCodeCheckStorageFailed 检查总存储大小失败
	ErrCodeCheckStorageFailed = 50001
	// ErrCodeCreateFileFailed 创建文件失败
	ErrCodeCreateFileFailed = 50002
	// ErrCodeSaveFileFailed 保存文件内容失败
	ErrCodeSaveFileFailed = 50003
	// ErrCodeAddMetadataFailed 添加文件元数据失败
	ErrCodeAddMetadataFailed = 50004
	// ErrCodeGetFilesFailed 获取文件列表失败
	ErrCodeGetFilesFailed = 50005
	// ErrCodeGetFileInfoFailed 获取文件信息失败
	ErrCodeGetFileInfoFailed = 50006
	// ErrCodeGetFileMetaFailed 获取文件信息失败（下载场景）
	ErrCodeGetFileMetaFailed = 50007
	// ErrCodeUpdateDownloadCountFailed 更新下载次数失败
	ErrCodeUpdateDownloadCountFailed = 50008
	// ErrCodeOpenFileFailed 打开文件失败
	ErrCodeOpenFileFailed = 50009
	// ErrCodeDeleteFileFailed 删除文件失败
	ErrCodeDeleteFileFailed = 50010
)
