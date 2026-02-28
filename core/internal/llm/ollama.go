package llm

import (
	"context"
	"fmt"
	"log"

	"github.com/Hirogava/WindowsAgent/core/internal/config"
	"github.com/Hirogava/WindowsAgent/core/internal/models"
	ollama "github.com/liliang-cn/ollama-go"
)

func SendTextToLLM(text string, llmURL string, prompt models.Prompts) (string, error) {
	ctx := context.Background()

	model, err := config.LoadOllamaConfigFromFile()
	if err != nil {
		return "", err
	}

	messages := []ollama.Message{
		{
			Role:    "system",
			Content: string(prompt),
		},
		{
			Role:    "user",
			Content: text,
		},
	}

	contextLLM := 2048

	options := ollama.Options{
		NumCtx: &contextLLM,
	}
	
	response, err := ollama.Chat(
        ctx,
        model,
        messages,
        func(r *ollama.ChatRequest) {
            r.Options = &options
        },
    )
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Message.Content)

	return response.Message.Content, nil
}
