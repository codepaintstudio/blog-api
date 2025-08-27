package services

import (
	"blog-api/internal/models"
	"blog-api/internal/repositories"
	"blog-api/internal/utils"
	"fmt"
	"mime/multipart"
	"gorm.io/gorm"
)

// FileService 文件服务接口
type FileService interface {
	Upload(userID uint, file *multipart.FileHeader) (*models.File, error)
	GetByID(id uint) (*models.File, error)
	Delete(id uint, userID uint) error
	GetUserFiles(userID uint, page, size int) ([]*models.File, int64, error)
	IncrementRefCount(id uint) error
	DecrementRefCount(id uint) error
	CleanupUnreferencedFiles() error
}

// fileService 文件服务实现
type fileService struct {
	fileRepo  repositories.FileRepository
	fileUtils *utils.FileUtils
}

// NewFileService 创建文件服务实例
func NewFileService(fileRepo repositories.FileRepository, fileUtils *utils.FileUtils) FileService {
	return &fileService{
		fileRepo:  fileRepo,
		fileUtils: fileUtils,
	}
}

// NewFileServiceDefault 创建默认文件服务实例（适配现有模式）
func NewFileServiceDefault() FileService {
	fileRepo := repositories.NewFileRepository()
	// 使用默认配置创建FileUtils
	fileUtils := utils.NewFileUtils(
		"./uploads",
		"http://localhost:8080/uploads",
		10485760, // 10MB
		[]string{"image/jpeg", "image/png", "image/gif", "image/webp"},
	)
	return &fileService{
		fileRepo:  fileRepo,
		fileUtils: fileUtils,
	}
}

// Upload 上传文件
func (s *fileService) Upload(userID uint, file *multipart.FileHeader) (*models.File, error) {
	// 1. 验证文件
	if err := s.fileUtils.ValidateFile(file); err != nil {
		return nil, fmt.Errorf("文件验证失败: %w", err)
	}

	// 2. 计算文件哈希值
	hash, err := s.fileUtils.CalculateFileHash(file)
	if err != nil {
		return nil, fmt.Errorf("计算文件哈希失败: %w", err)
	}

	// 3. 检查是否已存在相同文件
	existingFile, err := s.fileRepo.GetByHash(hash)
	if err == nil && existingFile != nil {
		// 文件已存在，增加引用计数
		existingFile.RefCount++
		if err := s.fileRepo.Update(existingFile); err != nil {
			return nil, fmt.Errorf("更新文件引用计数失败: %w", err)
		}
		return existingFile, nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}

	// 4. 生成文件名并保存文件
	filename := s.fileUtils.GenerateFilename(file.Filename, hash)
	filePath, err := s.fileUtils.SaveFile(file, filename)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 5. 生成文件URL
	fileURL := s.fileUtils.GenerateFileURL(filename)

	// 6. 创建文件记录
	fileModel := &models.File{
		UserID:       userID,
		Filename:     filename,
		OriginalName: file.Filename,
		Size:         file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		Path:         filePath,
		URL:          fileURL,
		Hash:         hash,
		RefCount:     1,
	}

	if err := s.fileRepo.Create(fileModel); err != nil {
		// 如果数据库操作失败，删除已保存的文件
		s.fileUtils.DeleteFile(filePath)
		return nil, fmt.Errorf("创建文件记录失败: %w", err)
	}

	return fileModel, nil
}

// GetByID 根据ID获取文件
func (s *fileService) GetByID(id uint) (*models.File, error) {
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文件不存在")
		}
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}
	return file, nil
}

// Delete 删除文件（减少引用计数）
func (s *fileService) Delete(id uint, userID uint) error {
	// 1. 获取文件信息
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("文件不存在")
		}
		return fmt.Errorf("查询文件失败: %w", err)
	}

	// 2. 检查权限（只有上传者可以删除）
	if file.UserID != userID {
		return fmt.Errorf("无权限删除此文件")
	}

	// 3. 减少引用计数
	file.RefCount--
	if file.RefCount <= 0 {
		// 引用计数为0，删除文件记录和物理文件
		if err := s.fileRepo.Delete(id); err != nil {
			return fmt.Errorf("删除文件记录失败: %w", err)
		}
		// 删除物理文件
		if err := s.fileUtils.DeleteFile(file.Path); err != nil {
			// 记录日志，但不影响主流程
			fmt.Printf("Warning: 删除物理文件失败: %s, error: %v\n", file.Path, err)
		}
	} else {
		// 更新引用计数
		if err := s.fileRepo.Update(file); err != nil {
			return fmt.Errorf("更新文件引用计数失败: %w", err)
		}
	}

	return nil
}

// GetUserFiles 获取用户文件列表
func (s *fileService) GetUserFiles(userID uint, page, size int) ([]*models.File, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	files, total, err := s.fileRepo.GetUserFiles(userID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户文件失败: %w", err)
	}

	return files, total, nil
}

// IncrementRefCount 增加文件引用计数
func (s *fileService) IncrementRefCount(id uint) error {
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return err
	}

	file.RefCount++
	return s.fileRepo.Update(file)
}

// DecrementRefCount 减少文件引用计数
func (s *fileService) DecrementRefCount(id uint) error {
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return err
	}

	file.RefCount--
	if file.RefCount <= 0 {
		// 引用计数为0，删除文件
		if err := s.fileRepo.Delete(id); err != nil {
			return err
		}
		// 删除物理文件
		if err := s.fileUtils.DeleteFile(file.Path); err != nil {
			fmt.Printf("Warning: 删除物理文件失败: %s, error: %v\n", file.Path, err)
		}
	} else {
		if err := s.fileRepo.Update(file); err != nil {
			return err
		}
	}

	return nil
}

// CleanupUnreferencedFiles 清理无引用的文件
func (s *fileService) CleanupUnreferencedFiles() error {
	// 获取所有无引用的文件
	files, err := s.fileRepo.GetUnreferencedFiles()
	if err != nil {
		return fmt.Errorf("查询无引用文件失败: %w", err)
	}

	var deletedCount int
	for _, file := range files {
		// 删除数据库记录
		if err := s.fileRepo.Delete(file.ID); err != nil {
			fmt.Printf("Warning: 删除文件记录失败: ID=%d, error: %v\n", file.ID, err)
			continue
		}

		// 删除物理文件
		if err := s.fileUtils.DeleteFile(file.Path); err != nil {
			fmt.Printf("Warning: 删除物理文件失败: %s, error: %v\n", file.Path, err)
		}

		deletedCount++
	}

	fmt.Printf("清理完成，删除了 %d 个无引用文件\n", deletedCount)
	return nil
}