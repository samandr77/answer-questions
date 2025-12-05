package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/andrey-samosuk/answer-questions/internal/api"
	"github.com/andrey-samosuk/answer-questions/internal/config"
	"github.com/andrey-samosuk/answer-questions/internal/repository"
	"github.com/andrey-samosuk/answer-questions/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка получения SQL DB: %v", err)
	}
	if err := runMigrations(sqlDB); err != nil {
		log.Fatalf("Ошибка при применении миграций: %v", err)
	}

	questionRepo := repository.NewQuestionRepository(db)
	answerRepo := repository.NewAnswerRepository(db)

	questionService := service.NewQuestionService(questionRepo)
	answerService := service.NewAnswerService(answerRepo, questionRepo)

	handler := api.NewHandler(questionService, answerService, cfg.Server.RequestTimeout)

	router := api.NewRouter(handler)
	mux := router.Setup()

	server := &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Запуск HTTP сервера на http://localhost:%s", cfg.Server.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка сервера: %v", err)
		}
	}()

	<-sigChan
	log.Println("Получен сигнал завершения, останавливаем приложение...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("Ошибка при закрытии БД: %v", err)
	}

	wg.Wait()
	log.Println("Приложение остановлено")
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	migrationsDir := "migrations"
	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}
