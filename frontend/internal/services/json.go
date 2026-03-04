package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type MicrophoneConfig struct {
	Device          string `json:"device"`
	DurationSeconds int    `json:"duration_seconds"`
	TriggerKey      string `json:"trigger_key"`
}

type ConfigService struct {
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (cs *ConfigService) LoadJsonFileNames() ([]string, error) {
	logger := getFrontendLogger()
	configDir, err := cs.configDir()
	if err != nil {
		logger.Printf("load json names failed: %v", err)
		return nil, err
	}

	dirs, err := os.ReadDir(configDir)
	if err != nil {
		logger.Printf("read config dir failed: %v", err)
		return nil, err
	}

	names := []string{}

	for _, dir := range dirs {
		if strings.HasSuffix(dir.Name(), ".json") {
			names = append(names, dir.Name())
		}
	}

	if len(names) == 0 {
		logger.Printf("no json configs found in: %s", configDir)
		return nil, nil
	}

	fmt.Println(names)
	logger.Printf("json config files loaded: %v", names)

	return names, nil
}

func (cs *ConfigService) ListMicrophones() ([]string, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-list_devices", "true",
		"-f", "dshow",
		"-i", "dummy",
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	_ = cmd.Run()

	lines := strings.Split(stderr.String(), "\n")
	devices := make([]string, 0)
	for _, line := range lines {
		if !strings.Contains(line, "audio") {
			continue
		}

		start := strings.Index(line, "\"")
		end := strings.LastIndex(line, "\"")
		if start == -1 || end == -1 || end <= start {
			continue
		}

		name := line[start+1 : end]
		if name != "" {
			devices = append(devices, name)
		}
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no microphones found")
	}

	return devices, nil
}

func (cs *ConfigService) LoadMicrophoneConfig() (*MicrophoneConfig, error) {
	logger := getFrontendLogger()
	path, err := cs.microphoneConfigPath()
	if err != nil {
		logger.Printf("microphone config path failed: %v", err)
		return nil, err
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		defaultCfg := &MicrophoneConfig{Device: "", DurationSeconds: 5, TriggerKey: "space"}
		if saveErr := cs.SaveMicrophoneConfig(defaultCfg); saveErr != nil {
			logger.Printf("save default microphone config failed: %v", saveErr)
			return nil, saveErr
		}
		logger.Printf("default microphone config created: %s", path)
		return defaultCfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Printf("read microphone config failed: %v", err)
		return nil, err
	}

	cfg := &MicrophoneConfig{}
	if err = json.Unmarshal(data, cfg); err != nil {
		logger.Printf("parse microphone config failed: %v", err)
		return nil, err
	}

	if cfg.DurationSeconds <= 0 {
		cfg.DurationSeconds = 5
	}

	if strings.TrimSpace(cfg.TriggerKey) == "" {
		cfg.TriggerKey = "space"
	}

	return cfg, nil
}

func (cs *ConfigService) SaveMicrophoneConfig(cfg *MicrophoneConfig) error {
	logger := getFrontendLogger()
	path, err := cs.microphoneConfigPath()
	if err != nil {
		logger.Printf("microphone config path failed: %v", err)
		return err
	}

	if cfg.DurationSeconds <= 0 {
		cfg.DurationSeconds = 5
	}

	if strings.TrimSpace(cfg.TriggerKey) == "" {
		cfg.TriggerKey = "space"
	}

	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		logger.Printf("marshal microphone config failed: %v", err)
		return err
	}

	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		logger.Printf("write microphone config failed: %v", err)
		return err
	}

	logger.Printf("microphone config saved: device=%s duration=%d key=%s", cfg.Device, cfg.DurationSeconds, cfg.TriggerKey)
	return nil
}

func (cs *ConfigService) ReadConfigFile(name string) ([]byte, error) {
	configDir, err := cs.configDir()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(filepath.Join(configDir, name))
}

func (cs *ConfigService) WriteConfigFile(name string, data []byte) error {
	configDir, err := cs.configDir()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(configDir, name), data, 0644)
}

func (cs *ConfigService) MapToJSONString(m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func (cs *ConfigService) configDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := cwd
	for {
		candidate := filepath.Join(current, "config")
		info, statErr := os.Stat(candidate)
		if statErr == nil && info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}

		current = parent
	}

	return "", fmt.Errorf("config directory not found from cwd: %s", cwd)
}

func (cs *ConfigService) repoRootDir() (string, error) {
	configDir, err := cs.configDir()
	if err != nil {
		return "", err
	}

	return filepath.Dir(configDir), nil
}

func (cs *ConfigService) microphoneConfigPath() (string, error) {
	configDir, err := cs.configDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "microphone-config.json"), nil
}
