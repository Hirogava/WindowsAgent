package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendTextToAudio(text string, ttsURL string) error {
	body := &bytes.Buffer{}

	req, err := http.NewRequest("GET", ttsURL, body)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("text", text)
	req.URL.RawQuery = q.Encode()

	fmt.Println("📡 Отправка текста в TTS...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("📡 Получение аудио от TTS...")

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TTS API error: %s", resp.Status)
	}

	file, err := os.Create("output.wav")
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
