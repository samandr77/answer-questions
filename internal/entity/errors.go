package entity

type ErrorResponse struct {
	Error string `json:"error"`
}

type CustomError struct {
	Code    int
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}

var (
	ErrQuestionNotFound = CustomError{
		Code:    404,
		Message: "Вопрос не найден",
	}
	ErrInvalidQuestionText = CustomError{
		Code:    400,
		Message: "Текст вопроса не может быть пустым",
	}
	ErrQuestionAlreadyExists = CustomError{
		Code:    409,
		Message: "Вопрос с таким текстом уже существует",
	}

	ErrAnswerNotFound = CustomError{
		Code:    404,
		Message: "Ответ не найден",
	}
	ErrInvalidAnswerText = CustomError{
		Code:    400,
		Message: "Текст ответа не может быть пустым",
	}
	ErrInvalidUserID = CustomError{
		Code:    400,
		Message: "ID пользователя не может быть пустым",
	}

	ErrDatabaseConnection = CustomError{
		Code:    500,
		Message: "Ошибка подключения к базе данных",
	}
	ErrDatabaseQuery = CustomError{
		Code:    500,
		Message: "Ошибка при выполнении запроса к базе данных",
	}

	ErrValidationFailed = CustomError{
		Code:    400,
		Message: "Ошибка валидации данных",
	}
)
