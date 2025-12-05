package repository

import (
	"context"

	"github.com/andrey-samosuk/answer-questions/internal/entity"

	"gorm.io/gorm"
)

type answerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &answerRepository{db: db}
}

func (r *answerRepository) Create(ctx context.Context, answer *entity.Answer) (*entity.Answer, error) {
	if err := r.db.WithContext(ctx).Create(answer).Error; err != nil {
		return nil, err
	}
	return answer, nil
}

func (r *answerRepository) GetByID(ctx context.Context, id int) (*entity.Answer, error) {
	// TODO: реализовать
	return nil, nil
}

func (r *answerRepository) GetByQuestionID(ctx context.Context, questionID int) ([]entity.Answer, error) {
	var answers []entity.Answer
	if err := r.db.WithContext(ctx).Where("question_id = ?", questionID).Order("created_at DESC").Find(&answers).Error; err != nil {
		return nil, err
	}
	if answers == nil {
		return []entity.Answer{}, nil
	}
	return answers, nil
}

func (r *answerRepository) Delete(ctx context.Context, id int) error {
	// TODO: реализовать
	return nil
}
