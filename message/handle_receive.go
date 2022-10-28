package message

import (
	"fmt"
	"net/http"

	"encore.dev/beta/errs"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

type TranscribeEvent struct {
	MediaUrl string
	From     string
	To       string
}

var TranscriptionRequests = pubsub.NewTopic[*TranscribeEvent]("transcriptions", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

func writeError(w http.ResponseWriter, err error) {
	rlog.Error("failed", "err", err)
	errs.HTTPError(w, err)
}

// ReceiveMessage is the webhook endpoint that Twilio will
// call whenever a WhatsApp message is received
//
// encore:api public raw path=/message/receive
func ReceiveMessage(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		writeError(w, err)
		return
	}

	message, err := formToMessage(req.Form)
	if err != nil {
		writeError(w, err)
		return
	}

	rlog.Info("received message",
		"num_media", message.NumMedia,
		"from", message.From,
		"content_type", message.MediaContentType0,
	)
	// if message.NumMedia > 1 {
	// 	// @TODO investigate: message may have more than 1 content
	// 	// type; currently dealing only with Type0
	// }

	if message.MediaContentType0 == "audio/ogg" {
		event := &TranscribeEvent{
			MediaUrl: message.MediaUrl0,
			From:     message.From,
			To:       message.To,
		}
		if _, err := TranscriptionRequests.Publish(req.Context(), event); err != nil {
			writeError(w, err)
			return
		}

		fmt.Fprint(w, "🤖 Audio file received. I'll reply with the transcription when it's ready.")
		return
	}

	msg := fmt.Sprintf(
		"🤖 Hello, %v! Forward your audio messages to me, and I'll text you back a transcription.",
		message.ProfileName,
	)
	fmt.Fprint(w, msg)
}
