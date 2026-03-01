package service

import (
	"os/exec"
)

func (ar *ActionRegistry) OpenUrlInBrowser(query string) error {
	cmd := exec.Command("cmd", "/C", "start", query)
	return cmd.Run()
}
