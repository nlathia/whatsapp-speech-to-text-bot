package twilio

import (
	"fmt"
	"net/http"

	"encore.dev/beta/errs"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

const (
	audioReceived = "🤖 Audio file received. I'll reply with the transcription when it's ready."
	textReceived  = "🤖 Hello, %v! Forward your audio messages to me, and I'll text you back a transcription."
)

type TwilioMessage struct {
	From              string
	MediaContentType0 string
	MediaUrl          string
	ProfileName       string
	To                string
}

var TranscriptionRequests = pubsub.NewTopic[*TwilioMessage]("transcriptions", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

// ReceiveMessage is the webhook endpoint that Twilio will
// call whenever a WhatsApp message is received
//
// encore:api public raw path=/message/receive
func ReceiveMessage(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		errs.HTTPError(w, err)
		return
	}

	message := &TwilioMessage{
		From:              req.Form["From"][0],
		MediaContentType0: req.Form["MediaContentType0"][0],
		MediaUrl:          req.Form["MediaUrl0"][0],
		ProfileName:       req.Form["ProfileName"][0],
	}

	rlog.Info("received message", "type", message.MediaContentType0)
	if message.MediaContentType0 == "audio/ogg" {
		if message.MediaUrl == "" {
			return
		}
		if _, err := TranscriptionRequests.Publish(req.Context(), message); err != nil {
			errs.HTTPError(w, err)
			return
		}

		fmt.Fprint(w, audioReceived)
		return
	}

	msg := fmt.Sprintf(textReceived, message.ProfileName)
	fmt.Fprint(w, msg)
}
