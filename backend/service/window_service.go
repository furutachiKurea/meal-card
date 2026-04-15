package service

import (
	"backend/model"
	"backend/repository"
)

// WindowService 窗口业务逻辑
type WindowService struct {
	windowRepo *repository.WindowRepository
}

// NewWindowService 创建 WindowService 实例
func NewWindowService(windowRepo *repository.WindowRepository) *WindowService {
	return &WindowService{windowRepo: windowRepo}
}

// ListWindows 获取所有窗口
func (s *WindowService) ListWindows() ([]model.Window, error) {
	return s.windowRepo.ListWindows()
}

// CreateWindow 创建窗口
func (s *WindowService) CreateWindow(name string) (*model.Window, error) {
	if name == "" {
		return nil, newBizError(ErrCodeValidationError, "窗口名称不能为空")
	}
	w := &model.Window{Name: name}
	if err := s.windowRepo.CreateWindow(w); err != nil {
		return nil, err
	}
	return w, nil
}
