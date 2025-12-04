package repository

import (
	"context"

	"github.com/andrey-samosuk/answer-questions/internal/entity"
)

type answerRepository struct {
	// db *gorm.DB - добавится позже
}

func NewAnswerRepository() AnswerRepository {
	return &answerRepository{}
}

func (r *answerRepository) Create(ctx context.Context, answer *entity.Answer) error {
	// TODO: реализовать
	return nil
}

func (r *answerRepository) GetByID(ctx context.Context, id int) (*entity.Answer, error) {
	// TODO: реализовать
	return nil, nil
}

func (r *answerRepository) GetByQuestionID(ctx context.Context, questionID int) ([]entity.Answer, error) {
	// TODO: реализовать
	return nil, nil
}

func (r *answerRepository) Delete(ctx context.Context, id int) error {
	// TODO: реализовать
	return nil
}
