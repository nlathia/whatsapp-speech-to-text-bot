package twilio

import (
	"context"

	"encore.app/whisper"
)

type TranscriptionParams struct {
	MediaUrl string
}

type TranscriptionResponse struct {
	Text     string
	Language string
}

// Transcribe returns a transcription of a mediaUrl
//
//encore:api private method=POST path=/message/transcribe
func Transcribe(ctx context.Context, params TranscriptionParams) (*TranscriptionResponse, error) {
	rsp, err := whisper.Transcribe(ctx, &whisper.WhisperRequest{
		MediaUrl: params.MediaUrl,
	})
	if err != nil {
		return nil, err
	}

	return &TranscriptionResponse{
		Text:     rsp.Text,
		Language: rsp.Language,
	}, nil
}
