import whisper

from src.log import logger


def load_model(model_type: str = "base") -> whisper.Whisper:
    logger.info(f"Loading: '{model_type}' model.")
    model = whisper.load_model("base")
    logger.info(f"Finished loading: '{model_type}' model.")
    return model


def predict(model, uri: str) -> dict:
    logger.info(f"Starting transcription: {uri}")
    result = model.transcribe(uri)
    logger.info(f"Finished transcription: {uri}")
    logger.info(f"Result: {result}")
    return result
