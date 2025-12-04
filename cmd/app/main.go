package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/andrey-samosuk/answer-questions/internal/api"
	"github.com/andrey-samosuk/answer-questions/internal/repository"
	"github.com/andrey-samosuk/answer-questions/internal/service"
)

func main() {
	// TODO: Инициализировать конфигурацию из переменных окружения
	// TODO: Подключиться к БД (GORM + PostgreSQL)
	// TODO: Применить миграции БД

	// Инициализация репозиториев
	questionRepo := repository.NewQuestionRepository()
	answerRepo := repository.NewAnswerRepository()

	// Инициализация сервисов
	questionService := service.NewQuestionService(questionRepo)
	answerService := service.NewAnswerService(answerRepo, questionRepo)

	// Инициализация обработчиков HTTP
	handler := api.NewHandler(questionService, answerService)

	// Конфигурация роутера
	router := api.NewRouter(handler)
	mux := router.Setup()

	// Создание HTTP сервера
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	server := &http.Server{
		Addr:         ":" + httpPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Канал для сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// WaitGroup для синхронизации горутин
	var wg sync.WaitGroup

	// Запуск сервера в отдельной горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Запуск HTTP сервера на http://localhost:%s", httpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("Получен сигнал завершения, останавливаем приложение...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	// TODO: Закрыть соединение с БД

	// Ожидание завершения всех горутин
	wg.Wait()
	log.Println("Приложение остановлено")
}
