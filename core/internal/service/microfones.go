package service

import (
	"bytes"
	"os/exec"
	"strings"
)

func GetMicrophones() ([]string, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-list_devices", "true",
		"-f", "dshow",
		"-i", "dummy",
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
	}

	output := stderr.String()

	var devices []string
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.Contains(line, "dshow") && strings.Contains(line, "audio") || strings.Contains(line, "video") {
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && end > start {
				devices = append(devices, line[start+1:end])
			}
		}
	}

	return devices, nil
}
