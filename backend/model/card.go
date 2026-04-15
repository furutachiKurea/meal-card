package model

import "time"

// CardStatus 饭卡状态
type CardStatus string

const (
	CardStatusActive    CardStatus = "active"    // 正常使用
	CardStatusLost      CardStatus = "lost"      // 已挂失
	CardStatusCancelled CardStatus = "cancelled" // 已注销
)

// Card 饭卡
type Card struct {
	ID           uint       `gorm:"primaryKey" json:"id"` // 自动编号，即卡号
	CardHolderID uint       `gorm:"not null" json:"cardHolderId"` // 持卡人ID
	CardHolder   CardHolder `gorm:"foreignKey:CardHolderID" json:"cardHolder,omitempty"`
	Deposit      int64      `gorm:"not null;default:0" json:"deposit"` // 押金，单位：分
	Balance      int64      `gorm:"not null;default:0" json:"balance"` // 余额，单位：分
	Status       CardStatus `gorm:"not null;default:'active'" json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
