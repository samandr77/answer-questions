package api // nolint: unused

import (
	"net/http"

	"github.com/andrey-samosuk/answer-questions/internal/service"
)

type Handler struct {
	questionService service.QuestionService
	answerService   service.AnswerService
}

func NewHandler(questionService service.QuestionService, answerService service.AnswerService) *Handler {
	return &Handler{
		questionService: questionService,
		answerService:   answerService,
	}
}

func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
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
