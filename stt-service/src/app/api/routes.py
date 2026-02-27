from fastapi import APIRouter, UploadFile, File, HTTPException
from ..services.whisper import transcribe_audio

router = APIRouter()

@router.get("/health")
async def health_check():
    return {"status": "ok"}

@router.post("/transcribe")
async def transcribe_audio_endpoint(file: UploadFile = File(...)):
    if file.content_type not in (
        "audio/mpeg",
        "audio/wav",
        "audio/x-wav",
        "audio/mp3",
        "application/octet-stream",
    ):
        raise HTTPException(status_code=400, detail="Unsupported audio type")

    # нужно передавать читаемый поток, а не байты целиком, иначе не работает
    text = transcribe_audio(await file.read())
    
    return {
        "transcription": text,
    }