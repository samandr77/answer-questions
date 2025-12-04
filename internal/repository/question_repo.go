package repository // nolint: unused

import (
	"context"

	"github.com/andrey-samosuk/answer-questions/internal/entity"
)

type questionRepository struct {
	// db *gorm.DB - добавится позже
}

func NewQuestionRepository() QuestionRepository {
	return &questionRepository{}
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
	// TODO: реализовать
	return nil, nil
}

func (r *questionRepository) Delete(ctx context.Context, id int) error {
	// TODO: реализовать
	return nil
}
