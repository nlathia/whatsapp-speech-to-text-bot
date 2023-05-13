package twilio

import (
	"context"
	"fmt"
	"time"

	"encore.dev/pubsub"
	"encore.dev/rlog"
)

const (
	noAudioResult = "😳 I couldn't detect any speech in that audio."
	audioResult   = "💬 *Transcription Result*\n\n%s"
)

var _ = pubsub.NewSubscription(
	TranscriptionRequests,
	"start-transcription",
	pubsub.SubscriptionConfig[*TwilioMessage]{
		Handler:     handleTranscribe,
		AckDeadline: time.Second * 60 * 5,
	},
)

// handleTranscribe transcribes a message and sends
// the result back to the user
func handleTranscribe(ctx context.Context, event *TwilioMessage) error {
	rlog.Info("event received", "media_url", event.MediaUrl)
	if event.MediaUrl == "" {
		return nil
	}

	transcription, err := Transcribe(ctx, TranscriptionParams{
		MediaUrl: event.MediaUrl,
	})
	if err != nil {
		rlog.Error("transcription failed", "err", err)
		return err
	}

	return SendMessage(ctx, &MessageParams{
		To:   event.From,
		From: event.To,
		Body: formatResult(transcription.Text),
	})
}

func formatResult(text string) string {
	if text == "" {
		return noAudioResult
	}
	return fmt.Sprintf(audioResult, text)
}
