package message

import (
	"fmt"
	"net/http"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type Message struct {
	AccountSid        string
	ApiVersion        string
	Body              string
	Forwarded         bool
	From              string
	MediaContentType0 string
	MediaUrl0         string
	MessageSid        string
	NumMedia          int
	NumSegments       int
	ProfileName       string
	ReferralNumMedia  int
	SmsMessageSid     string
	SmsSid            string
	SmsStatus         string
	To                string
	WaId              string
}

func formToMessage(form map[string][]string) (*Message, error) {
	values := firstValues(form)

	numMedia, err := parseInt("NumMedia", values)
	if err != nil {
		return nil, err
	}
	numSegments, err := parseInt("NumSegments", values)
	if err != nil {
		return nil, err
	}
	referralNumMedia, err := parseInt("ReferralNumMedia", values)
	if err != nil {
		return nil, err
	}

	fwdStatus, err := parseBool("Forwarded", values)
	if err != nil {
		return nil, err
	}

	return &Message{
		AccountSid:        values["AccountSid"],
		ApiVersion:        values["ApiVersion"],
		Body:              values["Body"],
		Forwarded:         fwdStatus,
		From:              values["From"],
		MediaContentType0: values["MediaContentType0"],
		MediaUrl0:         values["MediaUrl0"],
		MessageSid:        values["MessageSid"],
		NumMedia:          numMedia,
		NumSegments:       numSegments,
		ProfileName:       values["ProfileName"],
		ReferralNumMedia:  referralNumMedia,
		SmsMessageSid:     values["SmsMessageSid"],
		SmsSid:            values["SmsSid"],
		SmsStatus:         values["SmsStatus"],
		To:                values["To"],
		WaId:              values["WaId"],
	}, nil
}

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

	// Message may have more than 1 content type?

	if message.MediaContentType0 == "audio/ogg" {
		fmt.Fprintf(w, "Hello audio!")
		// Call whisper service
		return
	}

	// Reply saying that the file isn't audio
	// fmt.Fprintf(w, "Hello, %v! ", message.ProfileName)
}
