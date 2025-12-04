package entity

import "errors"

var (
	// Question errors
	ErrQuestionNotFound    = errors.New("question not found")
	ErrInvalidQuestionText = errors.New("invalid question text")

	// Answer errors
	ErrAnswerNotFound    = errors.New("answer not found")
	ErrInvalidAnswerText = errors.New("invalid answer text")
	ErrInvalidUserID     = errors.New("invalid user id")

	// Database errors
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")

	// Validation errors
	ErrValidationFailed = errors.New("validation failed")
)
