# WindowsAgent

Голосовой ассистент для Windows: запись голоса, распознавание речи, генерация ответа через LLM, выполнение действия и озвучка результата.

## Что реализовано сейчас

- `core/cmd/jarvis` (Go): оркестратор пайплайна.
- `stt-service` (Python/FastAPI): `POST /api/transcribe`.
- `tts-service` (Python/FastAPI): `GET /api/text-to-speech`.
- `action-service` (Go/Gin):
	- `POST /api/command-execute` — выполнить действие (сейчас открытие URL/поиска в браузере);
	- `POST /api/play-audio` — проиграть WAV, полученный от TTS.
	- `GET /api/wait-for-a-key-press` — ожидание нажатия `Space`.
- `core/internal/llm`: интеграция с Ollama (`config/ollama-config.json`).

## Текущий пайплайн

1. `jarvis` получает список микрофонов через `ffmpeg -f dshow -list_devices true`.
2. После выбора устройства ждёт нажатие `Space` через `action-service /api/wait-for-a-key-press`, затем записывает 5 секунд аудио и отправляет в STT (`:8001`).
3. Текст отправляется в LLM дважды:
	 - для голосового ответа (`PromptForTaskVoice`);
	 - для JSON-команды выполнения (`PromptForTaskExecution`).
4. Голосовой ответ уходит в TTS (`:8002`), затем WAV пересылается в `action-service /api/play-audio` (`:8003`) для воспроизведения.
5. JSON-команда уходит в `action-service /api/command-execute` (`:8003`).

## Зависимости

- Windows 10/11.
- Go (версия из `go.mod`).
- Python 3.10+.
- `ffmpeg` в `PATH`.
- Ollama (локально поднятый сервер) + модель из `config/ollama-config.json`.
- `piper` в `PATH` и модель TTS из `config/tts-config.json`.

## Запуск (локально)

Нужно 4 терминала из корня репозитория.

### 1) STT service

```powershell
cd .\stt-service
pip install -r requirements.txt
python -m uvicorn src.app.main:app --host 127.0.0.1 --port 8001
```

### 2) TTS service

```powershell
cd .\tts-service
pip install -r requirements.txt
python -m uvicorn src.app.main:app --host 127.0.0.1 --port 8002
```

### 3) Action service

```powershell
cd .\action-service\cmd
go run main.go
```

### 4) Jarvis orchestrator

```powershell
cd .\core\cmd\jarvis
go run main.go
```

## Порты и эндпоинты

- STT (`127.0.0.1:8001`)
	- `GET /api/health`
	- `POST /api/transcribe` (`multipart/form-data`, поле `file`)
- TTS (`127.0.0.1:8002`)
	- `GET /api/health`
	- `GET /api/text-to-speech?text=...` → `audio/wav`
- Action (`127.0.0.1:8003`)
	- `POST /api/command-execute`
	- `POST /api/play-audio` (`multipart/form-data`, поле `audio`)
	- `GET /api/wait-for-a-key-press`

## Конфигурация

- `config/ollama-config.json`
	- `model`: имя модели Ollama (текущее значение в репозитории: `qwen2.5:1.5b`).
- `config/stt-config.json`
	- `model`, `device`, `compute_type` для `faster-whisper` (текущие значения: `medium`, `cuda`, `float16`).
- `config/tts-config.json`
	- `model`: путь к `.onnx` файлу модели Piper (относительно `tts-service/src/app/services`).

## Важно

- В `jarvis` URL сервисов пока захардкожены в коде (`8001`, `8002`, `8003`).
- `action-service` в `command-execute` сейчас принимает только команды `browser` и `search`.
- Для браузерной команды используется `args[0]` как URL/строка запуска.