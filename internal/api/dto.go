package api

import "time"

type CreateQuestionRequest struct {
	Text string `json:"text"`
}

type QuestionResponse struct {
	ID        int              `json:"id"`
	Text      string           `json:"text"`
	CreatedAt time.Time        `json:"created_at"`
	Answers   []AnswerResponse `json:"answers,omitempty"`
}

type QuestionsListResponse struct {
	Questions []QuestionResponse `json:"questions"`
	Total     int                `json:"total"`
}

type QuestionMinimalResponse struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type QuestionsMinimalListResponse struct {
	Questions []QuestionMinimalResponse `json:"questions"`
}

type CreateAnswerRequest struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

type AnswerResponse struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	UserID     string    `json:"user_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
