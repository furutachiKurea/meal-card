package model

// Window 窗口机
type Window struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"` // 窗口名称
}
