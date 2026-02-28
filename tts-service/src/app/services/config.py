import json
from pathlib import Path


def LoadConfig() -> str:
    current = Path(__file__).resolve()

    for parent in current.parents:
        config_path = parent / 'config' / 'tts-config.json'
        if config_path.exists():
            with open(config_path, 'r', encoding='utf-8') as f:
                data = json.load(f)
            return data["model"]

    raise FileNotFoundError("Не найден файл config/tts-config.json")