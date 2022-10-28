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

	log := rlog.With("num_media", message.NumMedia, "from", message.From)
	if message.NumMedia > 1 {
		log.Info("received multiple media")
		// @TODO investigate: message may have more than 1 content
		// type; currently dealing only with Type0
	}

	if message.MediaContentType0 == "audio/ogg" {
		event := &TranscribeEvent{
			MediaUrl: message.MediaUrl0,
		}

		if _, err := TranscriptionRequests.Publish(req.Context(), event); err != nil {
			writeError(w, err)
			return
		}

		fmt.Fprint(w, "🤖 Audio file received. I'll reply with the transcription when it's ready.")
		return
	}

	log.Info("received non-audio format", "format", message.MediaContentType0)
	msg := fmt.Sprintf(
		"🤖 Hello, %v! Forward your audio messages to me, and I'll text you back a transcription.",
		message.ProfileName,
	)
	fmt.Fprint(w, msg)

	// if transcription.Text == "" {
	// 	fmt.Fprint(w, "😳 I couldn't detect any speech in that audio.")
	// 	return
	// }

	// // @TODO investigate formatting results by segment
	// // instead of just dumping it out
	// fmt.Fprintf(w, "💬 *Transcription Result*\n\n%s", transcription.Text)
}
