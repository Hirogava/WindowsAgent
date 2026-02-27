# WindowsAgent

Голосовой ассистент для Windows в стиле «JARVIS».

Текущий фокус проекта — базовый голосовой конвейер:

- запись голоса с микрофона;
- STT (speech-to-text) — распознавание речи;
- TTS (text-to-speech) — синтез речи;
- запуск из Go-клиента (`jarvis`), который связывает сервисы между собой.

## Что уже реализовано

- `core/cmd/jarvis`: CLI-сценарий «нажми Enter → запиши 5 сек → отправь в STT → отправь текст в TTS».
- `stt-service`: FastAPI-сервис для транскрибации аудио (`POST /api/transcribe`).
- `tts-service`: FastAPI-сервис для генерации WAV из текста (`GET /api/text-to-speech?text=...`).

На текущем этапе это MVP для голосового взаимодействия; логика «понимания намерений» и полноценного управления системой еще в развитии.

## Архитектура (сейчас)

1. `jarvis` записывает звук через `ffmpeg` (Windows `dshow`).
2. Аудио отправляется в `stt-service` как `multipart/form-data` с полем `file`.
3. Полученная транскрипция отправляется в `tts-service`.
4. `tts-service` возвращает WAV, который сохраняется как `output.wav`.

## Требования

- Windows 10/11.
- Go (совместимая версия по `go.mod`).
- Python 3.10+.
- `ffmpeg` в `PATH`.
- Для `tts-service`: установленный `piper` и голосовые модели:
	- `tts-service/src/app/services/ru_RU-ruslan-medium.onnx`
	- `tts-service/src/app/services/ru_RU-ruslan-medium.onnx.json`

## Быстрый старт

Откройте 3 терминала в корне репозитория.

### 1) Запуск STT

```powershell
cd .\stt-service
pip install -r requirements.txt
python -m uvicorn src.app.main:app --host 127.0.0.1 --port 8001
```

Проверка:

```powershell
Invoke-WebRequest -UseBasicParsing http://127.0.0.1:8001/api/health
```

### 2) Запуск TTS

```powershell
cd .\tts-service
pip install -r requirements.txt
python -m uvicorn src.app.main:app --host 127.0.0.1 --port 8002
```

Проверка:

```powershell
Invoke-WebRequest -UseBasicParsing http://127.0.0.1:8002/api/health
```

### 3) Запуск Jarvis (Go)

```powershell
cd .\core\cmd\jarvis
go run main.go
```

После нажатия Enter будет запись 5 секунд, отправка в STT и затем в TTS.

## API (актуально)

### STT

- `POST /api/transcribe`
- `multipart/form-data`
- поле файла: `file`
- поддерживаемые content-type: `audio/mpeg`, `audio/wav`, `audio/x-wav`, `audio/mp3`, `application/octet-stream`
- ответ:

```json
{
	"transcription": "..."
}
```

### TTS

- `GET /api/text-to-speech?text=<urlencoded_text>`
- ответ: бинарный `audio/wav`

## Структура репозитория

```text
WindowsAgent/
├─ core/
│  ├─ cmd/jarvis/                  # Go-точка входа ассистента
│  ├─ internal/
│  │  ├─ service/                  # запись аудио, STT/TTS запросы
│  │  ├─ handlers/                 # заготовки обработчиков
│  │  ├─ intents/                  # заготовки intent-логики
│  │  └─ ...
│  └─ pkg/
├─ stt-service/                    # Python FastAPI STT сервис
│  ├─ src/app/api/routes.py
│  ├─ src/app/services/whisper.py
│  └─ requirements.txt
├─ tts-service/                    # Python FastAPI TTS сервис
│  ├─ src/app/api/routes.py
│  ├─ src/app/services/piper.py
│  └─ requirements.txt
├─ action-service/                 # блок действий/автоматизаций (в развитии)
├─ memory-service/                 # сервис памяти (в развитии)
├─ shared/                         # общие материалы
└─ plans/                          # заметки и планы
```

## Дальше по roadmap

- NLU/intent-слой между STT и action-service.
- Маршрутизация команд в действия ОС (приложения, файлы, браузер, автоматизация).
- Улучшение диалога и контекстной памяти.
- Единый docker/dev-оркестратор для всех сервисов.

## Примечания

- Сейчас в `jarvis` зашиты локальные URL:
	- STT: `http://127.0.0.1:8001/api/transcribe`
	- TTS: `http://127.0.0.1:8002/api/text-to-speech`
- Имя микрофона берется как строка устройства Windows (`dshow`) и должно совпадать с вашей системой.