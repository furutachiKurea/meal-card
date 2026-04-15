package repository

import (
	"backend/model"

	"gorm.io/gorm"
)

// WindowRepository 窗口的增删改查
type WindowRepository struct {
	db *gorm.DB
}

// NewWindowRepository 创建 WindowRepository 实例
func NewWindowRepository(db *gorm.DB) *WindowRepository {
	return &WindowRepository{db: db}
}

// ListWindows 查询所有窗口
func (r *WindowRepository) ListWindows() ([]model.Window, error) {
	var windows []model.Window
	err := r.db.Find(&windows).Error
	return windows, err
}

// CreateWindow 创建窗口
func (r *WindowRepository) CreateWindow(w *model.Window) error {
	return r.db.Create(w).Error
}

// FindWindowByID 根据 ID 查询窗口
func (r *WindowRepository) FindWindowByID(id uint) (*model.Window, error) {
	var w model.Window
	err := r.db.First(&w, id).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}
