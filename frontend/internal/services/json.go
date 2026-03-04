package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ConfigService struct {
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (cs *ConfigService) LoadJsonFileNames() ([]string, error) {
	configDir, err := cs.configDir()
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	names := []string{}

	for _, dir := range dirs {
		if strings.HasSuffix(dir.Name(), ".json") {
			names = append(names, dir.Name())
		}
	}

	if len(names) == 0 {
		return nil, nil
	}

	fmt.Println(names)

	return names, nil
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
