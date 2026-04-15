package model

import "time"

// DepositRecord 存款记录，每次充值产生一条
type DepositRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CardID    uint      `gorm:"not null;index" json:"cardId"`
	Card      Card      `gorm:"foreignKey:CardID" json:"card,omitempty"`
	Amount    int64     `gorm:"not null" json:"amount"` // 存款金额，单位：分
	CreatedAt time.Time `json:"createdAt"`
}
