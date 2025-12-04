package repository

import (
	"context"

	"github.com/andrey-samosuk/answer-questions/internal/entity"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *entity.Question) error

	GetByID(ctx context.Context, id int) (*entity.Question, error)

	GetAll(ctx context.Context) ([]entity.Question, error)

	Delete(ctx context.Context, id int) error
}

type AnswerRepository interface {
	Create(ctx context.Context, answer *entity.Answer) error

	GetByID(ctx context.Context, id int) (*entity.Answer, error)

	GetByQuestionID(ctx context.Context, questionID int) ([]entity.Answer, error)

	Delete(ctx context.Context, id int) error
}
