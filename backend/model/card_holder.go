package model

import "time"

// CardHolder 持卡人，记录办卡人的身份信息
type CardHolder struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`       // 姓名
	IDNumber  string    `gorm:"uniqueIndex;not null" json:"idNumber"` // 证件号
	CreatedAt time.Time `json:"createdAt"`
}
