package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andrey-samosuk/answer-questions/internal/entity"
)

type mockQuestionService struct {
	getAll         func(ctx context.Context) ([]entity.Question, error)
	createQuestion func(ctx context.Context, text string) (*entity.Question, error)
	getQuestion    func(ctx context.Context, id int) (*entity.Question, error)
	deleteQuestion func(ctx context.Context, id int) error
}

func (m *mockQuestionService) GetAllQuestions(ctx context.Context) ([]entity.Question, error) {
	if m.getAll != nil {
		return m.getAll(ctx)
	}
	return nil, nil
}

func (m *mockQuestionService) CreateQuestion(ctx context.Context, text string) (*entity.Question, error) {
	if m.createQuestion != nil {
		return m.createQuestion(ctx, text)
	}
	return nil, nil
}

func (m *mockQuestionService) GetQuestion(ctx context.Context, id int) (*entity.Question, error) {
	if m.getQuestion != nil {
		return m.getQuestion(ctx, id)
	}
	return nil, nil
}

func (m *mockQuestionService) DeleteQuestion(ctx context.Context, id int) error {
	if m.deleteQuestion != nil {
		return m.deleteQuestion(ctx, id)
	}
	return nil
}

type mockAnswerService struct {
	getAnswersByQuestion func(ctx context.Context, questionID int) ([]entity.Answer, error)
}

func (m *mockAnswerService) CreateAnswer(ctx context.Context, questionID int, userID string, text string) (*entity.Answer, error) {
	return nil, nil
}

func (m *mockAnswerService) GetAnswer(ctx context.Context, id int) (*entity.Answer, error) {
	return nil, nil
}

func (m *mockAnswerService) GetAnswersByQuestion(ctx context.Context, questionID int) ([]entity.Answer, error) {
	if m.getAnswersByQuestion != nil {
		return m.getAnswersByQuestion(ctx, questionID)
	}
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

	if response.Error != "Ошибка при выполнении запроса к базе данных" {
		t.Errorf("expected error 'Ошибка при выполнении запроса к базе данных', got '%s'", response.Error)
	}
}

func TestCreateQuestion_Success(t *testing.T) {
	createdTime := time.Now()
	mockQService := &mockQuestionService{
		createQuestion: func(ctx context.Context, text string) (*entity.Question, error) {
			return &entity.Question{
				ID:        1,
				Text:      text,
				CreatedAt: createdTime,
			}, nil
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	body := CreateQuestionRequest{Text: "What is Go?"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateQuestion(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response["id"] != float64(1) {
		t.Errorf("expected id 1, got %v", response["id"])
	}

	if response["text"] != "What is Go?" {
		t.Errorf("expected text 'What is Go?', got %v", response["text"])
	}
}

func TestCreateQuestion_EmptyText(t *testing.T) {
	mockQService := &mockQuestionService{
		createQuestion: func(ctx context.Context, text string) (*entity.Question, error) {
			return nil, entity.ErrInvalidQuestionText
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	body := CreateQuestionRequest{Text: ""}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateQuestion(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Error != "Текст вопроса не может быть пустым" {
		t.Errorf("expected error 'Текст вопроса не может быть пустым', got '%s'", response.Error)
	}
}

func TestCreateQuestion_WhitespaceText(t *testing.T) {
	mockQService := &mockQuestionService{
		createQuestion: func(ctx context.Context, text string) (*entity.Question, error) {
			return nil, entity.ErrInvalidQuestionText
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	body := CreateQuestionRequest{Text: "   "}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateQuestion(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateQuestion_DatabaseError(t *testing.T) {
	mockQService := &mockQuestionService{
		createQuestion: func(ctx context.Context, text string) (*entity.Question, error) {
			return nil, entity.ErrDatabaseQuery
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	body := CreateQuestionRequest{Text: "What is Go?"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateQuestion(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func createTestRequest(method, path string, body []byte) *http.Request {
	if body != nil {
		return httptest.NewRequest(method, path, bytes.NewReader(body))
	}
	return httptest.NewRequest(method, path, nil)
}

func TestGetQuestion_Success(t *testing.T) {
	createdTime := time.Now()
	mockQService := &mockQuestionService{
		getQuestion: func(ctx context.Context, id int) (*entity.Question, error) {
			if id == 1 {
				return &entity.Question{
					ID:        id,
					Text:      "What is Go?",
					CreatedAt: createdTime,
				}, nil
			}
			return nil, entity.ErrQuestionNotFound
		},
	}

	mockAService := &mockAnswerService{
		getAnswersByQuestion: func(ctx context.Context, questionID int) ([]entity.Answer, error) {
			if questionID == 1 {
				return []entity.Answer{
					{
						ID:         1,
						QuestionID: questionID,
						UserID:     "user1",
						Text:       "Go is a language",
						CreatedAt:  createdTime.Add(time.Hour),
					},
				}, nil
			}
			return []entity.Answer{}, nil
		},
	}

	handler := NewHandler(mockQService, mockAService, 5)

	req := createTestRequest(http.MethodGet, "/questions/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetQuestion(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response["id"] != float64(1) {
		t.Errorf("expected id 1, got %v", response["id"])
	}

	if response["text"] != "What is Go?" {
		t.Errorf("expected text 'What is Go?', got %v", response["text"])
	}

	answers, ok := response["answers"].([]interface{})
	if !ok {
		t.Errorf("answers should be array, got %T", response["answers"])
		return
	}

	if len(answers) != 1 {
		t.Errorf("expected 1 answer, got %d", len(answers))
	}
}

func TestGetQuestion_NoAnswers(t *testing.T) {
	createdTime := time.Now()
	mockQService := &mockQuestionService{
		getQuestion: func(ctx context.Context, id int) (*entity.Question, error) {
			if id == 2 {
				return &entity.Question{
					ID:        id,
					Text:      "What is Go?",
					CreatedAt: createdTime,
				}, nil
			}
			return nil, entity.ErrQuestionNotFound
		},
	}

	mockAService := &mockAnswerService{
		getAnswersByQuestion: func(ctx context.Context, questionID int) ([]entity.Answer, error) {
			return []entity.Answer{}, nil
		},
	}

	handler := NewHandler(mockQService, mockAService, 5)

	req := createTestRequest(http.MethodGet, "/questions/2", nil)
	req.SetPathValue("id", "2")
	w := httptest.NewRecorder()

	handler.GetQuestion(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	answers, ok := response["answers"].([]interface{})
	if !ok {
		t.Errorf("answers should be array, got %T", response["answers"])
		return
	}

	if len(answers) != 0 {
		t.Errorf("expected 0 answers, got %d", len(answers))
	}
}

func TestGetQuestion_NotFound(t *testing.T) {
	mockQService := &mockQuestionService{
		getQuestion: func(ctx context.Context, id int) (*entity.Question, error) {
			return nil, entity.ErrQuestionNotFound
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	req := createTestRequest(http.MethodGet, "/questions/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetQuestion(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Error != "Вопрос не найден" {
		t.Errorf("expected error 'Вопрос не найден', got '%s'", response.Error)
	}
}

func TestGetQuestion_InvalidID(t *testing.T) {
	handler := NewHandler(&mockQuestionService{}, &mockAnswerService{}, 5)

	req := createTestRequest(http.MethodGet, "/questions/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetQuestion(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteQuestion_Success(t *testing.T) {
	mockQService := &mockQuestionService{
		deleteQuestion: func(ctx context.Context, id int) error {
			if id == 1 {
				return nil
			}
			return entity.ErrQuestionNotFound
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	req := createTestRequest(http.MethodDelete, "/questions/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteQuestion(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestDeleteQuestion_NotFound(t *testing.T) {
	mockQService := &mockQuestionService{
		deleteQuestion: func(ctx context.Context, id int) error {
			return entity.ErrQuestionNotFound
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	req := createTestRequest(http.MethodDelete, "/questions/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.DeleteQuestion(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Error != "Вопрос не найден" {
		t.Errorf("expected error 'Вопрос не найден', got '%s'", response.Error)
	}
}

func TestDeleteQuestion_InvalidID(t *testing.T) {
	handler := NewHandler(&mockQuestionService{}, &mockAnswerService{}, 5)

	req := createTestRequest(http.MethodDelete, "/questions/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.DeleteQuestion(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteQuestion_DatabaseError(t *testing.T) {
	mockQService := &mockQuestionService{
		deleteQuestion: func(ctx context.Context, id int) error {
			return entity.ErrDatabaseQuery
		},
	}

	handler := NewHandler(mockQService, &mockAnswerService{}, 5)

	req := createTestRequest(http.MethodDelete, "/questions/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteQuestion(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
