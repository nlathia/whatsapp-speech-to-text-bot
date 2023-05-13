package twilio

import (
	"fmt"
	"net/http"

	"encore.dev/beta/errs"
	"encore.dev/pubsub"
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
		From:              getFirst(req.Form, "From"),
		MediaContentType0: getFirst(req.Form, "MediaContentType0"),
		MediaUrl:          getFirst(req.Form, "MediaUrl0"),
		ProfileName:       getFirst(req.Form, "ProfileName"),
		To:                getFirst(req.Form, "To"),
	}
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

func getFirst(form map[string][]string, key string) string {
	values, exists := form[key]
	if !exists {
		return ""
	}
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
