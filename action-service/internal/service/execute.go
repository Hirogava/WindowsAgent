package service

import (
	"fmt"
	"os/exec"
)

func (ar *ActionRegistry) OpenUrlInBrowser(query string) error {
	cmd := exec.Command("cmd", "/C", "start", query)
	return cmd.Run()
}

func (ar *ActionRegistry) PlayWav(path string) error {
	cmd := exec.Command(
		"powershell",
		"-c",
		fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync();", path),
	)

	return cmd.Run()
}
