import whisper

model = whisper.load_model("small")

# load audio and pad/trim it to fit 30 seconds
audio = whisper.load_audio("73245.mp3")
audio = whisper.pad_or_trim(audio)

# make log-Mel spectrogram and move to the same device as the model
mel = whisper.log_mel_spectrogram(audio, n_mels=model.dims.n_mels).to(model.device)

# # detect the spoken language
# _, probs = model.detect_language(mel)
# print(f"Detected language: {max(probs, key=probs.get)}")

# decode the audio
options = whisper.DecodingOptions(
    language="ru",
    without_timestamps=True,
)
result = whisper.decode(model, mel, options)

# print the recognized text
print(result.text)