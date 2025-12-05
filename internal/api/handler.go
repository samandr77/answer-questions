package api // nolint: unused

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/andrey-samosuk/answer-questions/internal/entity"
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
		if customErr, ok := err.(entity.CustomError); ok {
			if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error: customErr.Message,
			}); err != nil {
				log.Printf("Ошибка при отправке ответа: %v", err)
			}
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
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	var req CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
			Error: "Некорректный формат JSON",
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
		return
	}

	question, err := h.questionService.CreateQuestion(ctx, req.Text)
	if err != nil {
		if customErr, ok := err.(entity.CustomError); ok {
			log.Printf("Ошибка при создании вопроса: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(customErr.Code)
			if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error: customErr.Message,
			}); err != nil {
				log.Printf("Ошибка при отправке ответа: %v", err)
			}
			return
		}

		log.Printf("Ошибка при создании вопроса: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
			Error: entity.ErrDatabaseQuery.Message,
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         question.ID,
		"text":       question.Text,
		"created_at": question.CreatedAt,
	}); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}

func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка парсинга ID: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
			Error: "Некорректный формат ID",
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
		return
	}

	question, err := h.questionService.GetQuestion(ctx, id)
	if err != nil {
		if customErr, ok := err.(entity.CustomError); ok {
			log.Printf("Ошибка при получении вопроса: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(customErr.Code)
			if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error: customErr.Message,
			}); err != nil {
				log.Printf("Ошибка при отправке ответа: %v", err)
			}
			return
		}

		log.Printf("Ошибка при получении вопроса: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
			Error: entity.ErrDatabaseQuery.Message,
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
		return
	}

	answers, err := h.answerService.GetAnswersByQuestion(ctx, id)
	if err != nil {
		log.Printf("Ошибка при получении ответов: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if customErr, ok := err.(entity.CustomError); ok {
			if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error: customErr.Message,
			}); err != nil {
				log.Printf("Ошибка при отправке ответа: %v", err)
			}
		}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         question.ID,
		"text":       question.Text,
		"created_at": question.CreatedAt,
		"answers":    answerResponses,
	}); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout)*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка парсинга ID: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
			Error: "Некорректный формат ID",
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
		return
	}

	err = h.questionService.DeleteQuestion(ctx, id)
	if err != nil {
		if customErr, ok := err.(entity.CustomError); ok {
			log.Printf("Ошибка при удалении вопроса: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(customErr.Code)
			if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error: customErr.Message,
			}); err != nil {
				log.Printf("Ошибка при отправке ответа: %v", err)
			}
			return
		}

		log.Printf("Ошибка при удалении вопроса: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
			Error: entity.ErrDatabaseQuery.Message,
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
			Error: "Некорректный формат ID",
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
		return
	}

	var req CreateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
			Error: "Некорректный формат JSON",
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
		return
	}

	answer, err := h.answerService.CreateAnswer(ctx, questionID, req.UserID, req.Text)
	if err != nil {
		if customErr, ok := err.(entity.CustomError); ok {
			log.Printf("Ошибка при создании ответа: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(customErr.Code)
			if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error: customErr.Message,
			}); err != nil {
				log.Printf("Ошибка при отправке ответа: %v", err)
			}
			return
		}

		log.Printf("Ошибка при создании ответа: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(entity.ErrorResponse{
			Error: entity.ErrDatabaseQuery.Message,
		}); err != nil {
			log.Printf("Ошибка при отправке ответа: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          answer.ID,
		"question_id": answer.QuestionID,
		"user_id":     answer.UserID,
		"text":        answer.Text,
		"created_at":  answer.CreatedAt,
	}); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}

func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("сервер запущен")); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}
