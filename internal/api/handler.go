package api // nolint: unused

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/andrey-samosuk/answer-questions/internal/service"
)

type Handler struct {
	questionService service.QuestionService
	answerService   service.AnswerService
	requestTimeout  int
}

func NewHandler(questionService service.QuestionService, answerService service.AnswerService, requestTimeout int) *Handler {
	return &Handler{
		questionService: questionService,
		answerService:   answerService,
		requestTimeout:  requestTimeout,
	}
}

func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	questions, err := h.questionService.GetAllQuestions(ctx)
	if err != nil {
		log.Printf("Ошибка при получении вопросов: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch questions",
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
		return
	}

	minimalQuestions := make([]QuestionMinimalResponse, len(questions))
	for i, q := range questions {
		minimalQuestions[i] = QuestionMinimalResponse{
			ID:   q.ID,
			Text: q.Text,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(QuestionsMinimalListResponse{
		Questions: minimalQuestions,
	}); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}
