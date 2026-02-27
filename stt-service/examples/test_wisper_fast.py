from faster_whisper import WhisperModel

model = WhisperModel(
    "medium",
    device="cuda",          # ВАЖНО
    compute_type="float16"  # Для RTX
)

segments, info = model.transcribe(
    "test.wav",
    vad_filter=True,
    beam_size=1,
    language="ru",
)

segments = list(segments)

print("Detected language '%s' with probability %f" % (info.language, info.language_probability))

for segment in segments:
    print("[%.2fs -> %.2fs] %s" % (segment.start, segment.end, segment.text))