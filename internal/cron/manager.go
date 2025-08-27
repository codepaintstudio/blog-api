package cron

import (
	"blog-api/internal/services"
	"fmt"
	"log"
	"time"
)

// Manager 定时任务管理器
type Manager struct {
	fileService services.FileService
	stopChan    chan bool
}

// NewManager 创建定时任务管理器
func NewManager(fileService services.FileService) *Manager {
	return &Manager{
		fileService: fileService,
		stopChan:    make(chan bool),
	}
}

// Start 启动定时任务
func (m *Manager) Start() {
	log.Println("定时任务管理器已启动")
	
	// 启动文件清理任务（每天凌晨2点执行）
	go m.startFileCleanupJob()
}

// Stop 停止定时任务
func (m *Manager) Stop() {
	log.Println("正在停止定时任务管理器...")
	m.stopChan <- true
	log.Println("定时任务管理器已停止")
}

// startFileCleanupJob 启动文件清理任务
func (m *Manager) startFileCleanupJob() {
	ticker := time.NewTicker(24 * time.Hour) // 每24小时检查一次
	defer ticker.Stop()

	// 立即执行一次清理（如果现在是凌晨2点左右）
	now := time.Now()
	if now.Hour() == 2 {
		m.executeFileCleanup()
	}

	// 计算下次凌晨2点的时间
	nextRun := m.getNextCleanupTime()
	log.Printf("下次文件清理任务将在 %s 执行", nextRun.Format("2006-01-02 15:04:05"))

	// 等待到下次执行时间
	time.AfterFunc(time.Until(nextRun), func() {
		m.executeFileCleanup()
		
		// 之后每24小时执行一次
		for {
			select {
			case <-ticker.C:
				m.executeFileCleanup()
			case <-m.stopChan:
				return
			}
		}
	})
}

// getNextCleanupTime 获取下次清理时间（明天凌晨2点）
func (m *Manager) getNextCleanupTime() time.Time {
	now := time.Now()
	
	// 今天凌晨2点
	today2AM := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	
	// 如果现在已经过了今天凌晨2点，则安排明天凌晨2点
	if now.After(today2AM) {
		return today2AM.Add(24 * time.Hour)
	}
	
	// 否则就是今天凌晨2点
	return today2AM
}

// executeFileCleanup 执行文件清理任务
func (m *Manager) executeFileCleanup() {
	log.Println("开始执行文件清理任务...")
	startTime := time.Now()

	if err := m.fileService.CleanupUnreferencedFiles(); err != nil {
		log.Printf("文件清理任务执行失败: %v", err)
	} else {
		duration := time.Since(startTime)
		log.Printf("文件清理任务执行完成，耗时: %v", duration)
	}
}

// CleanupTask 文件清理任务结构体
type CleanupTask struct {
	ExecuteAt time.Time `json:"execute_at"`
	Status    string    `json:"status"` // pending, running, completed, failed
	Message   string    `json:"message"`
	Duration  string    `json:"duration,omitempty"`
}

// GetCleanupTaskStatus 获取清理任务状态
func (m *Manager) GetCleanupTaskStatus() *CleanupTask {
	nextRun := m.getNextCleanupTime()
	
	return &CleanupTask{
		ExecuteAt: nextRun,
		Status:    "pending",
		Message:   fmt.Sprintf("下次清理任务将在 %s 执行", nextRun.Format("2006-01-02 15:04:05")),
	}
}

// ManualCleanup 手动触发清理任务
func (m *Manager) ManualCleanup() error {
	log.Println("手动触发文件清理任务...")
	
	startTime := time.Now()
	err := m.fileService.CleanupUnreferencedFiles()
	duration := time.Since(startTime)
	
	if err != nil {
		log.Printf("手动文件清理失败: %v", err)
		return err
	}
	
	log.Printf("手动文件清理完成，耗时: %v", duration)
	return nil
}