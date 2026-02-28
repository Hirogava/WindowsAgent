import json
from pathlib import Path


def LoadConfig():
    current = Path(__file__).resolve()

    for parent in current.parents:
        config_path = parent / 'config' / 'stt-config.json'
        if config_path.exists():
            with open(config_path, 'r', encoding='utf-8') as f:
                data = json.load(f)
            return data

    raise FileNotFoundError("Не найден файл config/stt-config.json")