package service

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
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
	return WaitForKeyPress("space")
}

func WaitForKeyPress(key string) error {
	if key == "" {
		key = "space"
	}

	fmt.Printf("Ожидаю нажатия клавиши %s...\n", key)
	uri := "http://127.0.0.1:8003/api/wait-for-a-key-press?key=" + url.QueryEscape(key)

	req, err := http.NewRequest("GET", uri, nil)
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
