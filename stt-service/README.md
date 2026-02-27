# STT Service

## Run

```powershell
python -m uvicorn --app-dir "WindowsAgent\stt-service" src.app.main:app --host 127.0.0.1 --port 8001
```

## Health check

```powershell
Invoke-WebRequest -UseBasicParsing http://127.0.0.1:8001/api/health
```

## Transcribe audio (curl)

```powershell
curl.exe -X POST "http://127.0.0.1:8001/api/transcribe" -F "file=@c:WindowsAgent\stt-service\examples\73245.mp3;type=audio/mpeg"
```

## Transcribe audio (PowerShell)

```powershell
Set-Location "WindowsAgent\stt-service"
Invoke-RestMethod -Uri "http://127.0.0.1:8001/api/transcribe" -Method Post -Form @{ file = Get-Item "examples\73245.mp3" }
```

## Notes

- Endpoint: `POST /api/transcribe`
- Form field name: `file`
- Supported upload content types: `audio/mpeg`, `audio/wav`, `audio/x-wav`, `audio/mp3`, `application/octet-stream`
