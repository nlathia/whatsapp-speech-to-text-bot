package transcribe

import (
	"context"
	"fmt"
	"time"

	"encore.app/message"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

var _ = pubsub.NewSubscription(
	message.TranscriptionRequests,
	"start-transcription",
	pubsub.SubscriptionConfig[*message.TranscribeEvent]{
		Handler:     StartTranscription,
		AckDeadline: time.Second * 60 * 5,
	},
)

func StartTranscription(ctx context.Context, event *message.TranscribeEvent) error {
	rlog.Info("event received", "media_url", event.MediaUrl)

	transcription, err := transcribe(ctx, event.MediaUrl)
	if err != nil {
		rlog.Error("transcription failed", "err", err)
		return err
	}

	var body string
	if transcription.Text == "" {
		body = "😳 I couldn't detect any speech in that audio."
	} else {
		body = fmt.Sprintf("💬 *Transcription Result*\n\n%s", transcription.Text)
	}

	return message.SendMessage(ctx, &message.MessageParams{
		To:   event.From,
		From: event.To,
		Body: body,
	})
}
