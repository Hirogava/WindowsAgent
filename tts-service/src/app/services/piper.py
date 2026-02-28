import subprocess
from pathlib import Path
from ..services.config import LoadConfig

def text_to_speech(text: str) -> bytes:
    input_file = Path(__file__).parent / "input.txt"
    output_file = Path(__file__).parent / "output.wav"

    with open(input_file, "w", encoding="utf-8") as f:
        f.write(text)
    
    model = LoadConfig()

    base_command = [
        "piper",
        "--model", str(Path(__file__).parent / model),
        "--input_file", str(input_file),
        "--output_file", str(output_file),
    ]

    try:
        subprocess.run(base_command + ["--cuda"], check=True)
    except subprocess.CalledProcessError:
        subprocess.run(base_command, check=True)
    
    with open(output_file, "rb") as f:
        file_data = f.read()

    output_file.unlink()
    input_file.unlink()
    return file_data