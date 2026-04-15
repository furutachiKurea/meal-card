package model

import "time"

// Transaction 消费记录，每次就餐扣款产生一条
type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CardID    uint      `gorm:"not null;index" json:"cardId"`
	Card      Card      `gorm:"foreignKey:CardID" json:"card,omitempty"`
	WindowID  uint      `gorm:"not null;index" json:"windowId"`
	Window    Window    `gorm:"foreignKey:WindowID" json:"window,omitempty"`
	Amount    int64     `gorm:"not null" json:"amount"` // 消费金额，单位：分
	CreatedAt time.Time `json:"createdAt"`
}
