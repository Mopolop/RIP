package main

import (
	"fmt"

	"db-integration/internal/app/config"
	"db-integration/internal/app/dsn"
	"db-integration/internal/app/handler"
	"db-integration/internal/app/repository"
	"db-integration/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// main инициализирует конфигурацию, репозиторий, хендлеры и запускает приложение
func main() {
	router := gin.Default() // создание нового роутера Gin

	// Загрузка конфигурации приложения
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	// Получение строки подключения к PostgreSQL
	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	// Инициализация репозитория
	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	// Создание хендлера с подключённым репозиторием
	hand := handler.NewHandler(rep)

	// Инициализация приложения и запуск сервера
	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
