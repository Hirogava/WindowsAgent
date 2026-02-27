package llm

import (
	"fmt"
	"log"
	"context"

	ollama "github.com/liliang-cn/ollama-go"
)

func SendTextToLLM(text string, llmURL string) (string, error) {
	ctx := context.Background()

	messages := []ollama.Message{
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
        "gemma3:4b",
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
