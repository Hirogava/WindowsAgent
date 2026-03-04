# core

Go-оркестратор голосового пайплайна (`jarvis`).

## Что делает

`core/cmd/jarvis/main.go` выполняет последовательность:

1. Получает список микрофонов через `ffmpeg` (`dshow`).
2. Читает `config/microphone-config.json` (микрофон, клавиша запуска и длительность записи).
3. Ждет нажатия выбранной клавиши через `http://127.0.0.1:8003/api/wait-for-a-key-press?key=...`.
4. Записывает аудио на длительность из конфига.
5. Отправляет аудио в STT (`http://127.0.0.1:8001/api/transcribe`).
6. Отправляет текст в Ollama:
   - запрос для озвучиваемого ответа;
   - запрос для JSON-команды выполнения.
7. Отправляет текст в TTS (`http://127.0.0.1:8002/api/text-to-speech`),
   затем пересылает WAV в action-service (`/api/play-audio`).
8. Отправляет JSON-команду в action-service (`/api/command-execute`).

## Рекомендуемый запуск

Сервис обычно запускается из frontend-кнопки `Запустить сервисы`.

```powershell
cd .\frontend\cmd\app
go run main.go
```

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

Используются:

- `config/ollama-config.json`
- `config/microphone-config.json`

Пример `microphone-config.json`:

```json
{
  "device": "Микрофон (fifine Microphone)",
  "duration_seconds": 5,
  "trigger_key": "space"
}
```

Пример `ollama-config.json`:

```json
{
  "model": "qwen2.5:1.5b"
}
```

Поиск файла идет от текущей рабочей директории вверх по дереву до каталога с `config/`.

## Логи

- `logs/jarvis.log` — основной runtime-лог оркестратора.

## Полезные детали

- Список устройств берет `GetMicrophones()` из вывода `ffmpeg -list_devices true`.
- Адреса сервисов в `main.go` пока захардкожены.
- Очистка ответа модели перед отправкой в action-service делается через `JsonCleaner()`.
