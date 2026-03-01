package service

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func SendTextToAudio(text string, ttsURL, actionServiceURL string) error {
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

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("audio", "output.wav")
	if err != nil {
		return err
	}

	_, err = io.Copy(part, resp.Body)
	if err != nil {
		return err
	}

	writer.Close()

	req, err = http.NewRequest("POST", actionServiceURL, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client = &http.Client{}
	_, err = client.Do(req)

	return nil
}
