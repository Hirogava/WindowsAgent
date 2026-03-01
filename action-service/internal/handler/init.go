package handler

import (
	"net/http"

	"github.com/Hirogava/WindowsAgent/action-service/internal/models"
	"github.com/Hirogava/WindowsAgent/action-service/internal/service"
	"github.com/gin-gonic/gin"
)

func InitHandlers(router *gin.Engine, ar service.ActionRegistry) {
	router.POST("/api/command-execute", func(ctx *gin.Context) {
		CommandExecute(ctx, ar)
	})
}

func CommandExecute(ctx *gin.Context, ar service.ActionRegistry) {
	var req models.PromptResponse

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Command != "browser" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "пока неизвестная команда"})
		return
	}

	if err := ar.OpenUrlInBrowser(req.Args[0]); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
