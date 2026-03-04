# WindowsAgent

Голосовой ассистент для Windows: запись голоса, распознавание речи, генерация ответа через LLM, выполнение действия и озвучка результата.

## Что реализовано сейчас

- `core/cmd/jarvis` (Go): оркестратор пайплайна.
- `frontend` (Go/Fyne): GUI для запуска/остановки сервисов, выбора микрофона и тестовой записи.
- `stt-service` (Python/FastAPI): `POST /api/transcribe`.
- `tts-service` (Python/FastAPI): `GET /api/text-to-speech`.
- `action-service` (Go/Gin):
  - `POST /api/command-execute` — выполнить действие (сейчас открытие URL/поиска в браузере).
  - `POST /api/play-audio` — проиграть WAV, полученный от TTS.
  - `GET /api/wait-for-a-key-press?key=space` — ожидание нажатия выбранной клавиши.
- `core/internal/llm`: интеграция с Ollama (`config/ollama-config.json`).

## Текущий пайплайн

1. В `frontend` выбирается микрофон, клавиша запуска записи и длительность, затем значения сохраняются в `config/microphone-config.json`.
2. `jarvis` получает список микрофонов через `ffmpeg -f dshow -list_devices true`.
3. `jarvis` читает `config/microphone-config.json`, ждёт нажатие выбранной клавиши через `action-service /api/wait-for-a-key-press?key=...`, записывает аудио и отправляет в STT (`:8001`).
4. Текст отправляется в LLM дважды:
   - для голосового ответа (`PromptForTaskVoice`);
   - для JSON-команды выполнения (`PromptForTaskExecution`).
5. Голосовой ответ уходит в TTS (`:8002`), затем WAV пересылается в `action-service /api/play-audio` (`:8003`) для воспроизведения.
6. JSON-команда уходит в `action-service /api/command-execute` (`:8003`).

## Зависимости

- Windows 10/11.
- Go (версия из `go.mod`).
- Python 3.10+.
- `ffmpeg` в `PATH`.
- Ollama (локально поднятый сервер) + модель из `config/ollama-config.json`.
- `piper` в `PATH` и модель TTS из `config/tts-config.json`.

## Запуск (локально)

Нужно 4 терминала из корня репозитория.

### Запуск через frontend (рекомендуется)

```powershell
cd .\frontend\cmd\app
go run main.go
```

В приложении на главной странице можно запускать/останавливать сервисы, а в меню `Настройки микрофона` выбрать устройство, клавишу записи и сделать тестовую запись.

## Документация модулей

- `core/README.md`
- `action-service/README.md`
- `stt-service/README.md`
- `tts-service/README.md`
- `config/README.md`

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
  - `GET /api/wait-for-a-key-press?key=space`

## Конфигурация

- `config/ollama-config.json`
  - `model`: имя модели Ollama (текущее значение в репозитории: `qwen2.5:1.5b`).
- `config/stt-config.json`
  - `model`, `device`, `compute_type` для `faster-whisper` (текущие значения: `medium`, `cuda`, `float16`).
- `config/tts-config.json`
  - `model`: путь к `.onnx` файлу модели Piper (относительно `tts-service/src/app/services`).
- `config/microphone-config.json`
  - `device`: выбранное имя микрофона из `ffmpeg`.
  - `duration_seconds`: длительность записи для `jarvis` и тестовой записи из frontend.
  - `trigger_key`: клавиша, которую ждёт `action-service` для запуска записи (`space`, `enter`, `a` и т.д.).

## Логи

- `logs/frontend.log` — действия GUI (запуск/остановка сервисов, сохранение конфига, тестовая запись).
- `logs/jarvis.log` — логи оркестратора `jarvis`.

## Важно

- В `jarvis` URL сервисов пока захардкожены в коде (`8001`, `8002`, `8003`).
- `action-service` в `command-execute` сейчас принимает только команды `browser` и `search`.
- Для браузерной команды используется `args[0]` как URL/строка запуска.
