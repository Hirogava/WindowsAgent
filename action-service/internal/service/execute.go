package service

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	hook "github.com/robotn/gohook"
)

func (ar *ActionRegistry) OpenUrlInBrowser(querys []string) error {
	query := ""
	standardGoogleSearch := "https://www.google.com/search?q="

	if len(querys) == 0 {
		fmt.Println("нет аргументов")
		return fmt.Errorf("нет аргументов для открытия в браузере")
	} else if len(querys) == 1 && strings.HasPrefix(querys[0], "http") {
		fmt.Println("по стандарту ", querys[0])
		query = querys[0]
	} else if strings.HasPrefix(querys[0], "http") {
		query = strings.Join(querys, "")
		query = strings.ReplaceAll(query, " ", "+")
		fmt.Println("составной ", query)
	} else {
		query = standardGoogleSearch + strings.Join(querys, "+")
		query = strings.ReplaceAll(query, " ", "+")
		fmt.Println("составной ", query)
	}

	cmd := exec.Command("cmd", "/C", "start", query)
	fmt.Println("комманда", cmd)
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

func (ar *ActionRegistry) ShutdownPC(args []string) error {
	time.Sleep(time.Second * 20)

	cmdArgs := append([]string{"/C", "shutdown"}, args...)
	cmd := exec.Command("cmd", cmdArgs...)
	return cmd.Run()
}

func (ar *ActionRegistry) RebootPC(args []string) error {
	time.Sleep(time.Second * 20)

	cmdArgs := append([]string{"/C", "shutdown"}, args...)
	cmd := exec.Command("cmd", cmdArgs...)
	return cmd.Run()
}

func (ar *ActionRegistry) WaitForKeyPress() error {
	fmt.Println("Ожидаю нажатия клавиши пробел...")

	for {
		if hook.AddEvent("space") {
			fmt.Println("Клавиша пробел нажата!")
			return nil
		}
	}
}
