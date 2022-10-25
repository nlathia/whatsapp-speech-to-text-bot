import logging
import sys


def get_logger():
    log = logging.getLogger(name="openai-transcribe")
    formatter = logging.Formatter("%(asctime)s - %(message)s")
    handler = logging.StreamHandler(sys.stderr)
    handler.setFormatter(formatter)
    log.setLevel(logging.INFO)
    log.addHandler(handler)
    return log


logger = get_logger()
