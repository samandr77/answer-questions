package entity

import "time"

type Answer struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	QuestionID int       `json:"question_id"`
	UserID     string    `json:"user_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	Question   *Question `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Answer) TableName() string {
	return "answers"
}
