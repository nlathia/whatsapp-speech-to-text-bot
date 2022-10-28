package message

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// ReceiveMessage is the webhook endpoint that Twilio will
// call whenever a WhatsApp message is received
//
// encore:api public raw path=/message/receive
func ReceiveMessage(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		rlog.Error("failed to parse form: %v", err)
		errs.HTTPError(w, err)
		return
	}

	message, err := formToMessage(req.Form)
	if err != nil {
		rlog.Error("failed to decode form: %v", err)
		errs.HTTPError(w, err)
		return
	}

	// @TODO investigate: message may have more than 1 content
	// type; currently dealing only with Type0
	if message.NumMedia > 1 {
		rlog.Info("received multiple media", "count", message.NumMedia)
	}

	// @TODO Check for other audio types (audio/mp3?)
	// image/jpeg
	if message.MediaContentType0 != "audio/ogg" {
		// @TODO write a better response, e.g.
		// reply saying that the file isn't audio
		rlog.Info("received non-audio format", "format", message.MediaContentType0)
		msg := fmt.Sprintf(
			"🤖 Hello, %v! Forward your audio messages to me, and I'll text you back a transcription!",
			message.ProfileName,
		)
		fmt.Fprint(w, msg)
		return
	}

	// @TODO this can be a very long-running call (for
	// large audio files). Move to running this async
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*60*2)
	defer func() {
		fmt.Print(w, "Timed-out while waiting for transcription")
		cancel()
	}()

	transcription, err := transcribe(ctx, message.MediaUrl0)
	if err != nil {
		rlog.Error("failed to transcribe")
		errs.HTTPError(w, err)
		return
	}

	// @TODO investigate formatting results by segment
	// instead of just dumping it out
	fmt.Fprintf(w, "💬 *Transcription Result*\n\n%s", transcription.Text)
}
