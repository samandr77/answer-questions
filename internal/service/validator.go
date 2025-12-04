package service

import "github.com/andrey-samosuk/answer-questions/internal/entity"

func ValidateQuestion(text string) error {
	// TODO: реализовать валидацию (пустота, длина)
	return nil
}

func ValidateAnswer(userID, text string) error {
	// TODO: реализовать валидацию (userID не пустой, текст не пустой)
	return nil
}

func ValidateQuestionNotNil(question *entity.Question) error {
	if question == nil {
		return entity.ErrQuestionNotFound
	}
	return nil
}
