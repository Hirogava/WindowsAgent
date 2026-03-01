package service

import (
	"context"

	"github.com/Hirogava/WindowsAgent/action-service/internal/models"
)

type Action interface {
	Name() string
	Validate(params map[string]interface{}) error
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

type ActionRegistry struct {
	Prompt models.PromptResponse
}

func NewActionRegistry() ActionRegistry {
	return ActionRegistry{}
}

func (ar *ActionRegistry) Name() string {
	return ar.Prompt.Command
}
