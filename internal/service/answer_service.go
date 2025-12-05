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
	if userID == "" {
		return nil, entity.ErrInvalidUserID
	}

	if text == "" {
		return nil, entity.ErrInvalidAnswerText
	}

	_, err := s.questionRepo.GetByID(ctx, questionID)
	if err != nil {
		return nil, entity.ErrQuestionNotFound
	}

	answer := &entity.Answer{
		QuestionID: questionID,
		UserID:     userID,
		Text:       text,
	}

	createdAnswer, err := s.answerRepo.Create(ctx, answer)
	if err != nil {
		return nil, entity.ErrDatabaseQuery
	}

	return createdAnswer, nil
}

func (s *answerService) GetAnswer(ctx context.Context, id int) (*entity.Answer, error) {
	answer, err := s.answerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entity.ErrAnswerNotFound
	}
	return answer, nil
}

func (s *answerService) GetAnswersByQuestion(ctx context.Context, questionID int) ([]entity.Answer, error) {
	answers, err := s.answerRepo.GetByQuestionID(ctx, questionID)
	if err != nil {
		return nil, entity.ErrDatabaseQuery
	}
	return answers, nil
}

func (s *answerService) DeleteAnswer(ctx context.Context, id int) error {
	if err := s.answerRepo.Delete(ctx, id); err != nil {
		if err == entity.ErrAnswerNotFound {
			return entity.ErrAnswerNotFound
		}
		return entity.ErrDatabaseQuery
	}
	return nil
}
