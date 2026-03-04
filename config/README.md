# Конфигурация WindowsAgent

Файлы конфигурации для всех сервисов проекта.

## microphone-config.json

```json
{
  "device": "Микрофон (fifine Microphone)",
  "duration_seconds": 5,
  "trigger_key": "space"
}
```

**Параметры:**
- `device` — имя микрофона из списка `ffmpeg`.
- `duration_seconds` — длительность записи в секундах для `jarvis` и тестовой записи из `frontend`.
- `trigger_key` — клавиша, которую ждёт `action-service` перед стартом записи (`space`, `enter`, `a` и т.д.).

**Как менять:**
- Рекомендуется менять через `frontend` в меню `Настройки микрофона`.
- Файл также можно редактировать вручную.

## ollama-config.json

```json
{
  "model": "qwen2.5:1.5b"
}
```

**Параметры:**
- `model` — имя модели Ollama для генерации ответов

**Рекомендации:**
- Для более точных ответов можно использовать более мощные модели, например:
  - `qwen2.5:1.5b` — улучшенная точность при небольшом размере
  - `qwen2.5:3b` — баланс производительности и качества
  - `llama3.2:3b` — хорошая альтернатива
- Модель должна быть предварительно загружена через `ollama pull <model>`

## stt-config.json

```json
{
  "model": "medium",
  "device": "cpu",
  "compute_type": "int8"
}
```

**Параметры:**
- `model` — размер модели Whisper (`tiny`, `base`, `small`, `medium`, `large`)
- `device` — устройство для инференса (`cpu` или `cuda`)
- `compute_type` — тип вычислений

**Использование GPU:**

Для ускорения работы STT на NVIDIA GPU измените:
```json
{
  "model": "medium",
  "device": "cuda",
  "compute_type": "float16"
}
```

Требования для GPU:
- NVIDIA GPU с поддержкой CUDA
- Установленный CUDA Toolkit
- `pip install faster-whisper[gpu]`

Возможные варианты `compute_type` для CUDA:
- `float16` — рекомендуется для GPU
- `int8` — меньше памяти, немного ниже качество

## tts-config.json

```json
{
  "model": "ru_RU-ruslan-medium.onnx"
}
```

**Параметры:**
- `model` — имя файла модели Piper (должен находиться в `tts-service/src/app/services/`)

**GPU для TTS:**

Piper автоматически пробует запуститься с флагом `--cuda`. Если GPU недоступен, происходит fallback на CPU.

## Примечания

- Все конфиги загружаются динамически: сервисы ищут папку `config/` вверх по дереву директорий от своего расположения.
- После изменения конфигов необходим перезапуск соответствующих сервисов.
