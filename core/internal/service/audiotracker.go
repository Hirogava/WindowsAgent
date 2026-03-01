package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	// "os"
	"os/exec"

	"github.com/Hirogava/WindowsAgent/core/internal/models"
	// "github.com/go-audio/audio"
	// "github.com/go-audio/wav"
	// "github.com/gordonklaus/portaudio"
)

var recording bool

func StartRecording() { recording = true /* запустить цикл записи */ }
func StopRecording()  { recording = false /* выйти из цикла */ }

// func RecordWav(filename string, durationSeconds int) error {
// 	const sampleRate = 16000
// 	const channels = 1

// 	frames := sampleRate * durationSeconds
// 	buffer := make([]int16, frames)

//  Инициализация PortAudio
// 	if err := portaudio.Initialize(); err != nil {
// 		return err
// 	}
// 	defer portaudio.Terminate()

// 	stream, err := portaudio.OpenDefaultStream(
// 		channels,
// 		0,
// 		float64(sampleRate),
// 		len(buffer),
// 		&buffer,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer stream.Close()

// 	fmt.Println("🎙 Запись...")

// 	if err := stream.Start(); err != nil {
// 		return err
// 	}

// 	if err := stream.Read(); err != nil {
// 		return err
// 	}

// 	if err := stream.Stop(); err != nil {
// 		return err
// 	}

// 	fmt.Println("✅ Запись завершена")

//  Создаём WAV файл
// 	outFile, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer outFile.Close()

// 	encoder := wav.NewEncoder(
// 		outFile,
// 		sampleRate,
// 		16,
// 		channels,
// 		1, // PCM
// 	)

// 	intBuf := &audio.IntBuffer{
// 		Data:           make([]int, len(buffer)),
// 		Format:         &audio.Format{NumChannels: channels, SampleRate: sampleRate},
// 		SourceBitDepth: 16,
// 	}

// 	for i, v := range buffer {
// 		intBuf.Data[i] = int(v)
// 	}

// 	if err := encoder.Write(intBuf); err != nil {
// 		return err
// 	}

// 	return encoder.Close()
// }

func RecordWav(device string, duration int) error {
	cmd := exec.Command(
		"ffmpeg",
		"-f", "dshow",
		"-i", fmt.Sprintf("audio=%s", device),
		"-t", fmt.Sprintf("%d", duration),
		"-ac", "1",
		"-ar", "16000",
		"output.wav",
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run()
}

func RecordAndSend(device string, duration int, sttURL string) (*models.STTResponse, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-f", "dshow",
		"-i", fmt.Sprintf("audio=%s", device),
		"-t", fmt.Sprintf("%d", duration),
		"-ac", "1",
		"-ar", "16000",
		"-f", "wav",
		"-",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	fmt.Println("🎙 Запись...")

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	fmt.Println("📡 Отправка в STT...")

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	fileWriter, err := writer.CreateFormFile("file", "audio.wav")
	if err != nil {
		return nil, err
	}

	if _, err := fileWriter.Write(out.Bytes()); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, sttURL, &requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("✅ Отправлено")

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stt error: %s", resp.Status)
	}

	var result models.STTResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
