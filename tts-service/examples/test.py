import subprocess
from pathlib import Path

subprocess.run(
    [
        "piper",
        "--model", str(Path(__file__).parent / "ru_RU-ruslan-medium.onnx"),
        "--input_file", str(Path(__file__).parent / "text.txt"),
        "--output_file", str(Path(__file__).parent / "out.wav"),
        "--cuda",
    ],
    check=True
)
print("OK")