# core

Go-оркестратор голосового пайплайна (`jarvis`).

## Что делает

`core/cmd/jarvis/main.go` выполняет последовательность:

1. Получает список микрофонов через `ffmpeg` (`dshow`).
2. Записывает 5 секунд аудио.
3. Отправляет аудио в STT (`http://127.0.0.1:8001/api/transcribe`).
4. Отправляет текст в Ollama:
   - запрос для озвучиваемого ответа;
   - запрос для JSON-команды выполнения.
5. Отправляет текст в TTS (`http://127.0.0.1:8002/api/text-to-speech`),
   затем пересылает WAV в action-service (`/api/play-audio`).
6. Отправляет JSON-команду в action-service (`/api/command-execute`).

## Запуск

```powershell
cd .\core\cmd\jarvis
go run main.go
```

## Требования

- `ffmpeg` в `PATH`.
- Доступные сервисы:
  - STT: `127.0.0.1:8001`
  - TTS: `127.0.0.1:8002`
  - action-service: `127.0.0.1:8003`
- Запущенный Ollama и модель из `config/ollama-config.json`.

## Конфиг

Используется `config/ollama-config.json`:

```json
{
  "model": "gemma3:1b"
}
```

Поиск файла идет от текущей рабочей директории вверх по дереву до каталога с `config/`.

## Полезные детали

- Список устройств берет `GetMicrophones()` из вывода `ffmpeg -list_devices true`.
- Адреса сервисов в `main.go` пока захардкожены.
- Очистка ответа модели перед отправкой в action-service делается через `JsonCleaner()`.
