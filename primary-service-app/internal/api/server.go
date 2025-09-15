package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"primary-service-app/internal/app/handler"
	"primary-service-app/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	// Правильная настройка статических файлов
	r.Static("/static", "./resources")
	r.LoadHTMLGlob("templates/*")

	// Добавляем корневой маршрут
	r.GET("/", handler.GetOrders)
	r.GET("/materials", handler.GetOrders)
	r.GET("/order/:id", handler.GetOrder)

	r.Run() // listen and serve on 0.0.0.0:8080
	log.Println("Server down")
}
