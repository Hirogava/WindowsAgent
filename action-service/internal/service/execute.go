package service

import (
	"fmt"
	"os/exec"
	"strings"
)

func (ar *ActionRegistry) OpenUrlInBrowser(querys []string) error {
	query := ""
	standardGoogleSearch := "https://www.google.com/search?q="

	if len(querys) == 0 {
		return fmt.Errorf("нет аргументов для открытия в браузере")
	} else if len(querys) == 1 && strings.HasPrefix(querys[0], "http") {
		query = querys[0]
	} else {
		query = standardGoogleSearch + strings.Join(querys, "+")
	}

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
