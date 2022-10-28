package transcribe

import (
	"context"

	"encore.app/message"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

var _ = pubsub.NewSubscription(
	message.TranscriptionRequests,
	"start-transcription",
	pubsub.SubscriptionConfig[*message.TranscribeEvent]{
		Handler: StartTranscription,
	},
)

func StartTranscription(ctx context.Context, event *message.TranscribeEvent) error {
	rlog.Info("event received", "media_url", event.MediaUrl)

	// transcription, err := transcribe(ctx, message.MediaUrl0)
	// if err != nil {
	// 	rlog.Error("transcription failed", "err", err)
	// 	return err
	// }

	return message.SendMessage(ctx, &message.MessageParams{
		To:   event.From,
		From: event.To,
		Body: "Here is the reply!",
	})
}
