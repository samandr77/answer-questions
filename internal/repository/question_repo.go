package repository

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
	if err := r.db.WithContext(ctx).Create(question).Error; err != nil {
		return err
	}
	return nil
}

func (r *questionRepository) GetByID(ctx context.Context, id int) (*entity.Question, error) {
	var question entity.Question
	if err := r.db.WithContext(ctx).First(&question, id).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) GetAll(ctx context.Context) ([]entity.Question, error) {
	var questions []entity.Question
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *questionRepository) GetByText(ctx context.Context, text string) (*entity.Question, error) {
	var question entity.Question
	if err := r.db.WithContext(ctx).Where("text = ?", text).First(&question).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&entity.Question{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entity.ErrQuestionNotFound
	}
	return nil
}
