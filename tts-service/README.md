# tts-service

Python/FastAPI сервис синтеза речи на базе `piper`.

## Запуск

```powershell
cd .\tts-service
pip install -r requirements.txt
python -m uvicorn src.app.main:app --host 127.0.0.1 --port 8002
```

## API

### `GET /api/health`

```json
{
  "status": "ok"
}
```

### `GET /api/text-to-speech?text=...`

- обязательный query-параметр: `text`
- ответ: бинарный `audio/wav`

## Проверка

```powershell
Invoke-WebRequest -UseBasicParsing "http://127.0.0.1:8002/api/health"
```

```powershell
$out = ".\examples\tts_http_test.wav"
Invoke-WebRequest -UseBasicParsing "http://127.0.0.1:8002/api/text-to-speech?text=Привет%20мир" -OutFile $out
(Get-Item $out).Length
```

## Конфиг

Используется `config/tts-config.json`:

```json
{
  "model": "ru_RU-ruslan-medium.onnx"
}
```

Значение `model` резолвится относительно `src/app/services/`.

## Внутренняя логика

- Временные файлы `input.txt` и `output.wav` создаются в `src/app/services/`.
- Сначала запускается `piper ... --cuda`, при ошибке выполняется fallback на CPU (без `--cuda`).
- После чтения результата временные файлы удаляются.

## Важно

- `piper` должен быть доступен в `PATH`.
- Файлы модели `.onnx` и `.onnx.json` должны быть рядом в `src/app/services/`.
