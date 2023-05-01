package message

import (
	"context"
	"encoding/json"

	"encore.dev/rlog"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type MessageParams struct {
	Body string
	From string
	To   string
}

//encore:api private method=POST path=/message/send
func SendMessage(ctx context.Context, p *MessageParams) error {
	client := getTwilioClient()

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(p.To)
	params.SetFrom(p.From)
	params.SetBody(p.Body)

	/*
		To receive real-time status updates for outbound messages, you can choose to set a
		Status Callback URL. Twilio sends a request to this URL each time your message status
		changes to one of the following: queued, failed, sent, delivered, read.
	*/

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		rlog.Error("error sending message", "err", err)
		return err
	} else {
		_, err = json.Marshal(*resp)
		if err != nil {
			rlog.Debug("failed to marshal result")
		}
	}
	return nil
}
