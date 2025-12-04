package entity

import "time"

type Question struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
}

func (Question) TableName() string {
	return "questions"
}
