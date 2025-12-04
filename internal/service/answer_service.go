package service

import (
	"context"

	"github.com/andrey-samosuk/answer-questions/internal/entity"
	"github.com/andrey-samosuk/answer-questions/internal/repository"
)

type AnswerService interface {
	CreateAnswer(ctx context.Context, questionID int, userID, text string) (*entity.Answer, error)

	GetAnswer(ctx context.Context, id int) (*entity.Answer, error)

	GetAnswersByQuestion(ctx context.Context, questionID int) ([]entity.Answer, error)

	DeleteAnswer(ctx context.Context, id int) error
}

type answerService struct {
	answerRepo   repository.AnswerRepository
	questionRepo repository.QuestionRepository
}

func NewAnswerService(answerRepo repository.AnswerRepository, questionRepo repository.QuestionRepository) AnswerService {
	return &answerService{
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
	}
}

func (s *answerService) CreateAnswer(ctx context.Context, questionID int, userID, text string) (*entity.Answer, error) {
	// TODO: реализовать проверку вопроса и валидацию
	return nil, nil
}

func (s *answerService) GetAnswer(ctx context.Context, id int) (*entity.Answer, error) {
	// TODO: реализовать
	return nil, nil
}

func (s *answerService) GetAnswersByQuestion(ctx context.Context, questionID int) ([]entity.Answer, error) {
	// TODO: реализовать
	return nil, nil
}

func (s *answerService) DeleteAnswer(ctx context.Context, id int) error {
	// TODO: реализовать
	return nil
}
