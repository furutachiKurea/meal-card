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
	ID           uint       `gorm:"primaryKey" json:"-"`                                // 内部自增主键，仅用于 FK，不对外暴露
	CardNo       string     `gorm:"uniqueIndex;not null;size:16" json:"cardNo"`         // 16 位业务卡号，对外唯一标识
	CardHolderID uint       `gorm:"not null" json:"cardHolderId"`                       // 持卡人 ID
	CardHolder   CardHolder `gorm:"foreignKey:CardHolderID" json:"cardHolder,omitempty"`
	Deposit      int64      `gorm:"not null;default:0" json:"deposit"` // 押金，单位：分，固定 2000
	Balance      int64      `gorm:"not null;default:0" json:"balance"` // 余额，单位：分
	Status       CardStatus `gorm:"not null;default:'active'" json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
