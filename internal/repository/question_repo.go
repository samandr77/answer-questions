package repository // nolint: unused

import (
	"context"

	"github.com/andrey-samosuk/answer-questions/internal/entity"

	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(ctx context.Context, question *entity.Question) error {
	// TODO: реализовать
	return nil
}

func (r *questionRepository) GetByID(ctx context.Context, id int) (*entity.Question, error) {
	// TODO: реализовать
	return nil, nil
}

func (r *questionRepository) GetAll(ctx context.Context) ([]entity.Question, error) {
	var questions []entity.Question
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *questionRepository) Delete(ctx context.Context, id int) error {
	// TODO: реализовать
	return nil
}
