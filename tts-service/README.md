# TTS Service

## Run

```powershell
python -m uvicorn --app-dir "WindowsAgent\tts-service" src.app.main:app --host 127.0.0.1 --port 8003
```

## Health check

```powershell
Invoke-WebRequest -UseBasicParsing "http://127.0.0.1:8002/api/health"
```

## Text-to-speech (save WAV)

```powershell
$out = "WindowsAgent\tts-service\examples\tts_http_test.wav"
Invoke-WebRequest -UseBasicParsing "http://127.0.0.1:8002/api/text-to-speech?text=Привет%20мир" -OutFile $out
(Get-Item $out).Length
```

## Endpoint

- `GET /api/text-to-speech?text=<urlencoded_text>`
- Response: `audio/wav`

## Notes

- Service tries `piper` with `--cuda`; if unavailable, falls back to CPU automatically.
- Voice model files are expected in `src/app/services`:
  - `ru_RU-ruslan-medium.onnx`
  - `ru_RU-ruslan-medium.onnx.json`
