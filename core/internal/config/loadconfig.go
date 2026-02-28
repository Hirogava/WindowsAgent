package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hirogava/WindowsAgent/core/internal/models"
)

func LoadOllamaConfigFromFile() (string, error) {
	filePath, err := findConfigPath("ollama-config.json")
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var cfg models.OllamaConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return "", err
	}

	return cfg.Model, err
}

func findConfigPath(fileName string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := wd
	for {
		candidate := filepath.Join(current, "config", fileName)
		if _, statErr := os.Stat(candidate); statErr == nil {
			return candidate, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("не найден файл config/%s (поиск от %s вверх по дереву)", fileName, wd)
}
