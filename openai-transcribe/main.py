from typing import List
import os

from flask import Flask, jsonify, request
from dataclasses import dataclass
from dataclasses_json import dataclass_json

from src.model.model import predict, load_model

app = Flask(__name__)
model = load_model()

@dataclass_json
@dataclass
class TranscriptionRequest:
    media_uri: str


@dataclass_json
@dataclass
class Segment:
    id: int
    seek: int
    start: float
    end: float
    text: str
    temperature: float
    avg_logprob: float
    compression_ratio: float
    no_speech_prob: float

    # tokens: List[int]


@dataclass_json
@dataclass
class TranscriptionResponse:
    text: str
    language: str
    segments: List[Segment]


@app.route('/', methods=['POST'])
def main():
    req = TranscriptionRequest.from_dict(request.get_json())
    result = predict(model=model, input=req.media_uri)

    resp = TranscriptionResponse.from_dict(result)
    return jsonify(resp.to_dict())


if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=int(os.environ.get('PORT', 8080)))
