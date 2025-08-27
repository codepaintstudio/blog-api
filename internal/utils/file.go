package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileUtils 文件工具结构体
type FileUtils struct {
	UploadPath string
	BaseURL    string
	MaxSize    int64
	AllowedTypes []string
}

// NewFileUtils 创建文件工具实例
func NewFileUtils(uploadPath, baseURL string, maxSize int64, allowedTypes []string) *FileUtils {
	return &FileUtils{
		UploadPath:   uploadPath,
		BaseURL:      baseURL,
		MaxSize:      maxSize,
		AllowedTypes: allowedTypes,
	}
}

// ValidateFile 验证上传的文件
func (fu *FileUtils) ValidateFile(file *multipart.FileHeader) error {
	// 检查文件大小
	if file.Size > fu.MaxSize {
		return fmt.Errorf("文件大小超过限制，最大允许 %d 字节", fu.MaxSize)
	}

	// 检查文件类型
	fileType := file.Header.Get("Content-Type")
	if !fu.isAllowedType(fileType) {
		return fmt.Errorf("不支持的文件类型: %s", fileType)
	}

	return nil
}

// isAllowedType 检查文件类型是否允许
func (fu *FileUtils) isAllowedType(mimeType string) bool {
	for _, allowedType := range fu.AllowedTypes {
		if allowedType == mimeType {
			return true
		}
	}
	return false
}

// CalculateFileHash 计算文件的SHA256哈希值
func (fu *FileUtils) CalculateFileHash(file *multipart.FileHeader) (string, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 计算哈希值
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GenerateFilename 生成唯一的文件名
func (fu *FileUtils) GenerateFilename(originalName, hash string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Format("20060102150405")
	
	// 使用哈希的前8位 + 时间戳作为文件名
	filename := fmt.Sprintf("%s_%s%s", hash[:8], timestamp, ext)
	return filename
}

// SaveFile 保存文件到本地
func (fu *FileUtils) SaveFile(file *multipart.FileHeader, filename string) (string, error) {
	// 确保上传目录存在
	if err := fu.ensureUploadDir(); err != nil {
		return "", err
	}

	// 构造文件路径
	filePath := filepath.Join(fu.UploadPath, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePath, nil
}

// ensureUploadDir 确保上传目录存在
func (fu *FileUtils) ensureUploadDir() error {
	if _, err := os.Stat(fu.UploadPath); os.IsNotExist(err) {
		return os.MkdirAll(fu.UploadPath, 0755)
	}
	return nil
}

// GenerateFileURL 生成文件访问URL
func (fu *FileUtils) GenerateFileURL(filename string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(fu.BaseURL, "/"), filename)
}

// DeleteFile 删除本地文件
func (fu *FileUtils) DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在，视为成功
	}
	return os.Remove(filePath)
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// GetMimeTypeFromExtension 根据扩展名获取MIME类型
func GetMimeTypeFromExtension(ext string) string {
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".bmp":  "image/bmp",
		".svg":  "image/svg+xml",
	}
	
	return mimeTypes[strings.ToLower(ext)]
}

// IsImageFile 检查是否为图片文件
func IsImageFile(mimeType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/bmp",
		"image/svg+xml",
	}
	
	for _, imageType := range imageTypes {
		if imageType == mimeType {
			return true
		}
	}
	return false
}