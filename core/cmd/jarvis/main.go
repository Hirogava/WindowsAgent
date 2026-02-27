package main

import (
	"fmt"
	"log"

	"github.com/Hirogava/WindowsAgent/core/internal/service"
)

func main() {
	fmt.Println("Нажми Enter для записи 5 секунд...")
	fmt.Scanln()

	resp, err := service.RecordAndSend("Микрофон (fifine Microphone)", 5, "http://127.0.0.1:8001/api/transcribe")
	if err != nil {
		log.Fatal(err)
	}

	// вот тут будет функция для понимания контекста команды и отправки текста в TTS

	err = service.SendTextToAudio(resp.Transcription, "http://127.0.0.1:8002/api/text-to-speech")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Транскрибированный текст:", resp.Transcription)
}
