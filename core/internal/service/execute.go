package service

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendToActionService(jsonResp, url string) error {
	body := bytes.NewBuffer([]byte(jsonResp))

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("action service error: %s", resp.Status)
	}

	return nil
}

func WaitForSpaceKeyPress() error {
	fmt.Println("Ожидаю нажатия клавиши Space...")

	req, err := http.NewRequest("GET", "http://127.0.0.1:8003/api/wait-for-a-key-press", nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("action service error: %s", resp.Status)
	}

	return nil
}
