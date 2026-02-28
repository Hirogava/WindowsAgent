package models

type PromptResponse struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}
