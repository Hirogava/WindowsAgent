from io import BytesIO
from faster_whisper import WhisperModel

model = WhisperModel(
    "medium",
    device="cuda",
    compute_type="float16"
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