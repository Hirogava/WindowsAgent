package main

import (
	"fmt"
	"log"

	"github.com/Hirogava/WindowsAgent/core/internal/service"
)

func main() {
	microfones, err := service.GetMicrophones()
	if err != nil {
		log.Fatal(err)
	}

	if len(microfones) == 0 {
		log.Fatal("Нет доступных микрофонов")
	}

	fmt.Println("Доступные микрофоны:")
	for i, mic := range microfones {
		fmt.Printf("%d: %s\n", i+1, mic)
	}

	fmt.Println("Выберите микрофон для записи (введите номер):")
	var choice int
	_, err = fmt.Scanln(&choice)

	if err != nil || choice < 1 || choice > len(microfones) {
		log.Fatal("Неверный выбор микрофона")
	}

	mic := microfones[choice-1]

	fmt.Println("Нажми Enter для записи 5 секунд...")
	fmt.Scanln()

	resp, err := service.RecordAndSend(mic, 5, "http://127.0.0.1:8001/api/transcribe")
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
