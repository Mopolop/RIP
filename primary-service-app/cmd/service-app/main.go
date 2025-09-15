package main

import (
	"log"
	"primary-service-app/internal/api"
)

func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}
