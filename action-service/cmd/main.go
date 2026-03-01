package main

import (
	"log"

	"github.com/Hirogava/WindowsAgent/action-service/internal/handler"
	"github.com/Hirogava/WindowsAgent/action-service/internal/router"
	"github.com/Hirogava/WindowsAgent/action-service/internal/service"
)

func main() {
	ar := service.NewActionRegistry()

	r := router.NewRouter()
	log.Println("Add router")

	handler.InitHandlers(r, ar)
	log.Println("Init routes")

	log.Println("Server start")
	router.StartServer(r)
}
