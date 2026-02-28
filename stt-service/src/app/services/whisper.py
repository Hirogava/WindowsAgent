from io import BytesIO
from faster_whisper import WhisperModel
from ..services.config import LoadConfig

modelData = LoadConfig()

model = WhisperModel(
    modelData["model"],
    device=modelData["device"], # либо cpu
    compute_type=modelData["compute_type"] # int8
)

def transcribe_audio(data: bytes) -> str:
    segments, info = model.transcribe(
        BytesIO(data),
        vad_filter=True,
        beam_size=1,
        language="ru",
    )

    segments = list(segments)

    print("Detected language '%s' with probability %f" % (info.language, info.language_probability))

    result = ""

    for segment in segments:
        print("[%.2fs -> %.2fs] %s" % (segment.start, segment.end, segment.text))
        result += segment.text + " "
    
    return result.strip()