package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Hirogava/WindowsAgent/core/internal/llm"
	"github.com/Hirogava/WindowsAgent/core/internal/models"
	"github.com/Hirogava/WindowsAgent/core/internal/service"
)

type microphoneConfig struct {
	Device          string `json:"device"`
	DurationSeconds int    `json:"duration_seconds"`
	TriggerKey      string `json:"trigger_key"`
}

func main() {
	setupLogger()

	microfones, err := service.GetMicrophones()
	if err != nil {
		log.Fatal(err)
	}

	if len(microfones) == 0 {
		log.Fatal("Нет доступных микрофонов")
	}

	fmt.Println("Доступные микрофоны:")
	for i, mic := range microfones {
		fmt.Printf("%d: %s\n", i+1, mic)
	}

	var mic string
	cfg, err := loadMicrophoneFromConfig()
	if err != nil || cfg.Device == "" {
		mic = microfones[0]
		fmt.Printf("Микрофон из конфига не найден, используется: %s\n", mic)
	} else {
		mic = cfg.Device
		fmt.Printf("Выбран микрофон из конфига: %s\n", mic)
	}

	triggerKey := "space"
	recordDuration := 5
	if err == nil {
		if cfg.TriggerKey != "" {
			triggerKey = cfg.TriggerKey
		}
		if cfg.DurationSeconds > 0 {
			recordDuration = cfg.DurationSeconds
		}
	}

	fmt.Printf("Нажми %s для записи %d секунд...\n", triggerKey, recordDuration)

	for {
		updatedCfg, readErr := loadMicrophoneFromConfig()
		if readErr == nil {
			if updatedCfg.Device != "" {
				mic = updatedCfg.Device
			}
			if updatedCfg.TriggerKey != "" {
				triggerKey = updatedCfg.TriggerKey
			}
			if updatedCfg.DurationSeconds > 0 {
				recordDuration = updatedCfg.DurationSeconds
			}
		}

		err = service.WaitForKeyPress(triggerKey)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := service.RecordAndSend(mic, recordDuration, "http://127.0.0.1:8001/api/transcribe")
		if err != nil {
			log.Fatal(err)
		}

		responseVoice, err := llm.SendTextToLLM(resp.Transcription, "", models.PromptForTaskVoice)
		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			err = service.SendTextToAudio(responseVoice, "http://127.0.0.1:8002/api/text-to-speech", "http://127.0.0.1:8003/api/play-audio")
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()

		respJson, err := llm.SendTextToLLM(resp.Transcription, "", models.PromptForTaskExecution)
		if err != nil {
			log.Fatal(err)
		}

		respJson = service.JsonCleaner(respJson)

		err = service.SendToActionService(respJson, "http://127.0.0.1:8003/api/command-execute")

		fmt.Println("Ответ от LLM (команда для выполнения):", respJson)
		fmt.Println("Ответ от LLM для озвучки действий: ", responseVoice)
		fmt.Println("Транскрибированный текст:", resp.Transcription)
		wg.Wait()
	}
}

func setupLogger() {
	rootDir, err := repoRootDir()
	if err != nil {
		log.Printf("logger init warning: %v", err)
		return
	}

	logsDir := filepath.Join(rootDir, "logs")
	if err = os.MkdirAll(logsDir, 0755); err != nil {
		log.Printf("logger init warning: %v", err)
		return
	}

	logFile := filepath.Join(logsDir, "jarvis.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("logger init warning: %v", err)
		return
	}

	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.Printf("logging initialized: %s", logFile)
}

func repoRootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := cwd
	for {
		candidate := filepath.Join(current, "config")
		if info, statErr := os.Stat(candidate); statErr == nil && info.IsDir() {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("repo root not found from cwd: %s", cwd)
}

func loadMicrophoneFromConfig() (*microphoneConfig, error) {
	configPath, err := microphoneConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := &microphoneConfig{DurationSeconds: 5, TriggerKey: "space"}
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.DurationSeconds <= 0 {
		cfg.DurationSeconds = 5
	}

	if cfg.TriggerKey == "" {
		cfg.TriggerKey = "space"
	}

	return cfg, nil
}

func microphoneConfigPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := cwd
	for {
		candidate := filepath.Join(current, "config", "microphone-config.json")
		if _, statErr := os.Stat(candidate); statErr == nil {
			return candidate, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("microphone config not found from cwd: %s", cwd)
}
