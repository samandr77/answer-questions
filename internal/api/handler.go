package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
		sendCustomError(w, err, "Ошибка при получении вопросов")
		return
	}

	minimalQuestions := make([]QuestionMinimalResponse, len(questions))
	for i, q := range questions {
		minimalQuestions[i] = QuestionMinimalResponse{
			ID:   q.ID,
			Text: q.Text,
		}
	}

	sendJSON(w, http.StatusOK, QuestionsMinimalListResponse{
		Questions: minimalQuestions,
	})
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	var req CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		sendError(w, http.StatusBadRequest, "Некорректный формат JSON")
		return
	}

	question, err := h.questionService.CreateQuestion(ctx, req.Text)
	if err != nil {
		sendCustomError(w, err, "Ошибка при создании вопроса")
		return
	}

	sendJSON(w, http.StatusCreated, map[string]interface{}{
		"id":         question.ID,
		"text":       question.Text,
		"created_at": question.CreatedAt,
	})
}

func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка парсинга ID: %v", err)
		sendError(w, http.StatusBadRequest, "Некорректный формат ID")
		return
	}

	question, err := h.questionService.GetQuestion(ctx, id)
	if err != nil {
		sendCustomError(w, err, "Ошибка при получении вопроса")
		return
	}

	answers, err := h.answerService.GetAnswersByQuestion(ctx, id)
	if err != nil {
		sendCustomError(w, err, "Ошибка при получении ответов")
		return
	}

	answerResponses := make([]map[string]interface{}, len(answers))
	for i, a := range answers {
		answerResponses[i] = map[string]interface{}{
			"id":         a.ID,
			"user_id":    a.UserID,
			"text":       a.Text,
			"created_at": a.CreatedAt,
		}
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"id":         question.ID,
		"text":       question.Text,
		"created_at": question.CreatedAt,
		"answers":    answerResponses,
	})
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка парсинга ID: %v", err)
		sendError(w, http.StatusBadRequest, "Некорректный формат ID")
		return
	}

	err = h.questionService.DeleteQuestion(ctx, id)
	if err != nil {
		sendCustomError(w, err, "Ошибка при удалении вопроса")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	questionID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка парсинга ID: %v", err)
		sendError(w, http.StatusBadRequest, "Некорректный формат ID")
		return
	}

	var req CreateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		sendError(w, http.StatusBadRequest, "Некорректный формат JSON")
		return
	}

	answer, err := h.answerService.CreateAnswer(ctx, questionID, req.UserID, req.Text)
	if err != nil {
		sendCustomError(w, err, "Ошибка при создании ответа")
		return
	}

	sendJSON(w, http.StatusCreated, map[string]interface{}{
		"id":          answer.ID,
		"question_id": answer.QuestionID,
		"user_id":     answer.UserID,
		"text":        answer.Text,
		"created_at":  answer.CreatedAt,
	})
}

func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка парсинга ID: %v", err)
		sendError(w, http.StatusBadRequest, "Некорректный формат ID")
		return
	}

	answer, err := h.answerService.GetAnswer(ctx, id)
	if err != nil {
		sendCustomError(w, err, "Ошибка при получении ответа")
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"id":          answer.ID,
		"question_id": answer.QuestionID,
		"user_id":     answer.UserID,
		"text":        answer.Text,
		"created_at":  answer.CreatedAt,
	})
}

func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка парсинга ID: %v", err)
		sendError(w, http.StatusBadRequest, "Некорректный формат ID")
		return
	}

	err = h.answerService.DeleteAnswer(ctx, id)
	if err != nil {
		sendCustomError(w, err, "Ошибка при удалении ответа")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("сервер запущен")); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}
