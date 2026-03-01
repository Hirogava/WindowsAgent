package handler

import (
	"io"
	"os"
	"net/http"

	"github.com/Hirogava/WindowsAgent/action-service/internal/models"
	"github.com/Hirogava/WindowsAgent/action-service/internal/service"
	"github.com/gin-gonic/gin"
)

func InitHandlers(router *gin.Engine, ar service.ActionRegistry) {
	router.POST("/api/command-execute", func(ctx *gin.Context) {
		CommandExecute(ctx, ar)
	})

	router.POST("/api/play-audio", func(ctx *gin.Context) {
		GetAudioAndPlay(ctx, ar)
	})
}

func CommandExecute(ctx *gin.Context, ar service.ActionRegistry) {
	var req models.PromptResponse

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Command != "browser" && req.Command != "search" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "пока неизвестная команда"})
		return
	}

	if err := ar.OpenUrlInBrowser(req.Args); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func GetAudioAndPlay(ctx *gin.Context, ar service.ActionRegistry) {
	fileHeader, err := ctx.FormFile("audio")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "*.wav")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = ar.PlayWav(tmpFile.Name())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "played",
	})
}
