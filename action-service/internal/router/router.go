package router

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func StartServer(router *gin.Engine) {
	router.Run(":8003")
}
