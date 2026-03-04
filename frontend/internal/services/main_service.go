package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type MainService struct {
	mu        sync.Mutex
	processes map[string]*exec.Cmd
	config    *ConfigService
	logger    *log.Logger
}

func NewMainService() *MainService {
	return &MainService{
		processes: make(map[string]*exec.Cmd),
		config:    NewConfigService(),
		logger:    getFrontendLogger(),
	}
}

func (ms *MainService) StartAllServices() error {
	services := []struct {
		name string
		dir  string
		cmd  string
		args []string
	}{
		{name: "stt", dir: filepath.Join("stt-service"), cmd: "python", args: []string{"-m", "uvicorn", "src.app.main:app", "--host", "127.0.0.1", "--port", "8001"}},
		{name: "tts", dir: filepath.Join("tts-service"), cmd: "python", args: []string{"-m", "uvicorn", "src.app.main:app", "--host", "127.0.0.1", "--port", "8002"}},
		{name: "action", dir: filepath.Join("action-service", "cmd"), cmd: "go", args: []string{"run", "main.go"}},
		{name: "jarvis", dir: filepath.Join("core", "cmd", "jarvis"), cmd: "go", args: []string{"run", "main.go"}},
	}

	for _, service := range services {
		if err := ms.StartService(service.name, service.dir, service.cmd, service.args...); err != nil {
			ms.logger.Printf("start all services failed on %s: %v", service.name, err)
			return err
		}
	}

	ms.logger.Printf("all services started")

	return nil
}

func (ms *MainService) StopAllServices() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for name, cmd := range ms.processes {
		if cmd == nil || cmd.Process == nil {
			delete(ms.processes, name)
			continue
		}

		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		ms.logger.Printf("service stopped: %s", name)
		delete(ms.processes, name)
	}

	return nil
}

func (ms *MainService) RunningServicesCount() int {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	count := 0
	for _, cmd := range ms.processes {
		if cmd != nil && cmd.Process != nil {
			count++
		}
	}

	return count
}

func (ms *MainService) StartService(name string, relativeDir string, command string, args ...string) error {
	ms.mu.Lock()
	if existing, ok := ms.processes[name]; ok && existing != nil && existing.Process != nil {
		ms.mu.Unlock()
		ms.logger.Printf("service already running: %s", name)
		return nil
	}
	ms.mu.Unlock()

	rootDir, err := ms.config.repoRootDir()
	if err != nil {
		return err
	}

	workingDir := filepath.Join(rootDir, relativeDir)
	if _, err = os.Stat(workingDir); err != nil {
		return fmt.Errorf("directory not found for %s: %w", name, err)
	}

	cmd := exec.Command(command, args...)
	cmd.Dir = workingDir
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err = cmd.Start(); err != nil {
		ms.logger.Printf("service start failed: %s, err: %v", name, err)
		return fmt.Errorf("failed to start %s: %w", name, err)
	}

	ms.mu.Lock()
	ms.processes[name] = cmd
	ms.mu.Unlock()
	ms.logger.Printf("service started: %s, pid: %d", name, cmd.Process.Pid)

	go func() {
		_ = cmd.Wait()
		ms.mu.Lock()
		current, ok := ms.processes[name]
		if ok && current == cmd {
			ms.logger.Printf("service exited: %s", name)
			delete(ms.processes, name)
		}
		ms.mu.Unlock()
	}()

	return nil
}

func (ms *MainService) RecordAudioSample(duration int) (string, error) {
	if duration <= 0 {
		duration = 5
	}

	cfg, err := ms.config.LoadMicrophoneConfig()
	if err != nil {
		return "", err
	}

	if cfg.Device == "" {
		ms.logger.Printf("record sample failed: microphone not selected")
		return "", fmt.Errorf("microphone is not selected")
	}

	rootDir, err := ms.config.repoRootDir()
	if err != nil {
		return "", err
	}

	recordingsDir := filepath.Join(rootDir, "recordings")
	if err = os.MkdirAll(recordingsDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("sample-%d.wav", time.Now().Unix())
	fullPath := filepath.Join(recordingsDir, filename)

	cmd := exec.Command(
		"ffmpeg",
		"-f", "dshow",
		"-i", fmt.Sprintf("audio=%s", cfg.Device),
		"-t", fmt.Sprintf("%d", duration),
		"-ac", "1",
		"-ar", "16000",
		fullPath,
	)

	if err = cmd.Run(); err != nil {
		ms.logger.Printf("record sample failed: device=%s duration=%d err=%v", cfg.Device, duration, err)
		return "", err
	}

	ms.logger.Printf("record sample complete: %s", fullPath)

	return fullPath, nil
}
