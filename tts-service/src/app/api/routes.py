from fastapi import APIRouter, HTTPException, Response
from ..services.piper import text_to_speech

router = APIRouter()

@router.get("/health")
async def health_check():
    return {"status": "ok"}

@router.get("/text-to-speech")
async def text_to_speech_endpoint(text: str):
    if not text:
        raise HTTPException(status_code=400, detail="Text is required")
    
    audio_data = text_to_speech(text)

    return Response(content=audio_data, media_type="audio/wav")