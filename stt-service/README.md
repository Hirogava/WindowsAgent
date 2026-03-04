# stt-service

Python/FastAPI сервис распознавания речи на базе `faster-whisper`.

## Рекомендуемый запуск

Обычно сервис запускается из frontend-кнопки `Запустить сервисы`:

```powershell
cd .\frontend\cmd\app
go run main.go
```

## Запуск

```powershell
cd .\stt-service
pip install -r requirements.txt
python -m uvicorn src.app.main:app --host 127.0.0.1 --port 8001
```

## API

### `GET /api/health`

```json
{
	"status": "ok"
}
```

### `POST /api/transcribe`

- `multipart/form-data`
- поле: `file`
- поддерживаемые content type:
	- `audio/mpeg`
	- `audio/wav`
	- `audio/x-wav`
	- `audio/mp3`
	- `application/octet-stream`

Ответ:

```json
{
	"transcription": "..."
}
```

## Проверка

```powershell
Invoke-WebRequest -UseBasicParsing http://127.0.0.1:8001/api/health
```

```powershell
Invoke-RestMethod `
	-Uri "http://127.0.0.1:8001/api/transcribe" `
	-Method Post `
	-Form @{ file = Get-Item ".\examples\73245.mp3" }
```

## Конфиг

Используется файл `config/stt-config.json`:

```json
{
	"model": "medium",
	"device": "cpu",
	"compute_type": "int8"
}
```

Сервис ищет конфиг, поднимаясь от `src/app/services/config.py` вверх по директориям.

## Заметки

- Язык в `transcribe()` сейчас зафиксирован как `ru`.
- Включен `vad_filter=True`.
