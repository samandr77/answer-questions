package service

import (
	"context"

	"github.com/andrey-samosuk/answer-questions/internal/entity"
	"github.com/andrey-samosuk/answer-questions/internal/repository"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, text string) (*entity.Question, error)

	GetQuestion(ctx context.Context, id int) (*entity.Question, error)

	GetAllQuestions(ctx context.Context) ([]entity.Question, error)

	DeleteQuestion(ctx context.Context, id int) error
}

type questionService struct {
	repo repository.QuestionRepository
}

func NewQuestionService(repo repository.QuestionRepository) QuestionService {
	return &questionService{repo: repo}
}

func (s *questionService) CreateQuestion(ctx context.Context, text string) (*entity.Question, error) {
	if err := ValidateQuestion(text); err != nil {
		return nil, err
	}

	_, err := s.repo.GetByText(ctx, text)
	if err == nil {
		return nil, entity.ErrQuestionAlreadyExists
	}

	question := &entity.Question{Text: text}
	if err := s.repo.Create(ctx, question); err != nil {
		return nil, entity.ErrDatabaseQuery
	}

	return question, nil
}

func (s *questionService) GetQuestion(ctx context.Context, id int) (*entity.Question, error) {
	question, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, entity.ErrQuestionNotFound
	}

	return question, nil
}

func (s *questionService) GetAllQuestions(ctx context.Context) ([]entity.Question, error) {
	questions, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, entity.ErrDatabaseQuery
	}
	if questions == nil {
		return []entity.Question{}, nil
	}
	return questions, nil
}

func (s *questionService) DeleteQuestion(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if err == entity.ErrQuestionNotFound {
			return entity.ErrQuestionNotFound
		}
		return entity.ErrDatabaseQuery
	}
	return nil
}
