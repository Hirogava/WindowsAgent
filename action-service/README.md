# action-service

Go-сервис для выполнения действий и локального воспроизведения аудио.

## Что делает

- Поднимает HTTP API на `127.0.0.1:8003`.
- Принимает JSON-команду от `core` и выполняет действие.
- Принимает WAV-файл и воспроизводит его через Windows Media.SoundPlayer.
- Предоставляет endpoint ожидания нажатия выбранной клавиши для запуска записи в `core`.

## Рекомендуемый запуск

Сервис обычно запускается из frontend-кнопки `Запустить сервисы`.

```powershell
cd .\frontend\cmd\app
go run main.go
```

## Запуск

```powershell
cd .\action-service\cmd
go run main.go
```

## API

### Health

Отдельного health-эндпоинта пока нет.

### `GET /api/wait-for-a-key-press?key=space`

- блокирует запрос до нажатия указанной клавиши (`key`)
- используется `core` перед стартом записи с микрофона

Успешный ответ:

```json
{
  "status": "key_pressed",
  "key": "space"
}
```

### `POST /api/command-execute`

Входной JSON:

```json
{
  "command": "browser",
  "args": ["https://www.google.com/search?q=golang"]
}
```

Текущее поведение:

- Разрешены только `command`: `browser` или `search`.
- Для открытия используется `args[0]`.
- Команда исполняется через `cmd /C start <args[0]>`.

Успешный ответ:

```json
{
  "status": "success"
}
```

### `POST /api/play-audio`

- `multipart/form-data`
- поле файла: `audio`
- ожидается WAV

Успешный ответ:

```json
{
  "status": "played"
}
```

## Пример запроса (PowerShell)

```powershell
Invoke-RestMethod `
  -Uri "http://127.0.0.1:8003/api/command-execute" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"command":"browser","args":["https://ya.ru"]}'
```

## Важно

- Сервис рассчитан на Windows (`cmd`, `powershell` внутри исполнения).
- Логика валидации команд пока минимальная.
