import whisper

from src.log import logger


def load_model(model_type: str = "base.en") -> whisper.Whisper:
    logger.info(f"Loading: '{model_type}' model.")
    model = whisper.load_model("base")
    logger.info(f"Finished loading: '{model_type}' model.")
    return model


def predict(model, input: str) -> dict:
    logger.info(f"Starting transcription: {input}")
    result = model.transcribe(input)
    logger.info(f"Finished transcription: {input}")
    logger.info(f"Result: {result}")
    return result
