package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrey-samosuk/answer-questions/internal/entity"
)

type mockQuestionService struct {
	getAll func(ctx context.Context) ([]entity.Question, error)
}

func (m *mockQuestionService) GetAllQuestions(ctx context.Context) ([]entity.Question, error) {
	return m.getAll(ctx)
}

func (m *mockQuestionService) CreateQuestion(ctx context.Context, text string) (*entity.Question, error) {
	return nil, nil
}

func (m *mockQuestionService) GetQuestion(ctx context.Context, id int) (*entity.Question, error) {
	return nil, nil
}

func (m *mockQuestionService) DeleteQuestion(ctx context.Context, id int) error {
	return nil
}

type mockAnswerService struct{}

func (m *mockAnswerService) CreateAnswer(ctx context.Context, questionID int, userID string, text string) (*entity.Answer, error) {
	return nil, nil
}

func (m *mockAnswerService) GetAnswer(ctx context.Context, id int) (*entity.Answer, error) {
	return nil, nil
}

func (m *mockAnswerService) GetAnswersByQuestion(ctx context.Context, questionID int) ([]entity.Answer, error) {
	return nil, nil
}

func (m *mockAnswerService) DeleteAnswer(ctx context.Context, id int) error {
	return nil
}

func TestGetQuestions_Success(t *testing.T) {
	questions := []entity.Question{
		{
			ID:   1,
			Text: "What is Go?",
		},
		{
			ID:   2,
			Text: "How to use GORM?",
		},
	}

	mockQService := &mockQuestionService{
		getAll: func(ctx context.Context) ([]entity.Question, error) {
			return questions, nil
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	w := httptest.NewRecorder()

	handler.GetQuestions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var response QuestionsMinimalListResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(response.Questions) != 2 {
		t.Errorf("expected 2 questions, got %d", len(response.Questions))
	}

	if response.Questions[0].ID != 1 || response.Questions[0].Text != "What is Go?" {
		t.Errorf("expected first question with id=1 and text='What is Go?', got id=%d and text='%s'",
			response.Questions[0].ID, response.Questions[0].Text)
	}

	if response.Questions[1].ID != 2 || response.Questions[1].Text != "How to use GORM?" {
		t.Errorf("expected second question with id=2 and text='How to use GORM?', got id=%d and text='%s'",
			response.Questions[1].ID, response.Questions[1].Text)
	}
}

func TestGetQuestions_EmptyList(t *testing.T) {
	mockQService := &mockQuestionService{
		getAll: func(ctx context.Context) ([]entity.Question, error) {
			return []entity.Question{}, nil
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	w := httptest.NewRecorder()

	handler.GetQuestions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response QuestionsMinimalListResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(response.Questions) != 0 {
		t.Errorf("expected 0 questions, got %d", len(response.Questions))
	}
}

func TestGetQuestions_ServiceError(t *testing.T) {
	mockQService := &mockQuestionService{
		getAll: func(ctx context.Context) ([]entity.Question, error) {
			return nil, entity.ErrDatabaseQuery
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	w := httptest.NewRecorder()

	handler.GetQuestions(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Error != "database_error" {
		t.Errorf("expected error 'database_error', got '%s'", response.Error)
	}
}
