package twilio

import (
	"context"
	"encoding/json"

	"encore.dev/rlog"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type MessageParams struct {
	Body string
	From string
	To   string
}

var secrets struct {
	TwilioAccountSid string
	TwilioAuthToken  string
}

// SendMessage sends a message via Twilio
//
//encore:api private method=POST path=/message/send
func SendMessage(ctx context.Context, p *MessageParams) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: secrets.TwilioAccountSid,
		Password: secrets.TwilioAuthToken,
	})

	resp, err := client.Api.CreateMessage(&twilioApi.CreateMessageParams{
		To:   &p.To,
		From: &p.From,
		Body: &p.Body,
	})
	if err != nil {
		rlog.Error("error sending message", "err", err)
		return err
	}
	_, err = json.Marshal(*resp)
	if err != nil {
		rlog.Debug("failed to marshal result")
		// Ignore error
	}
	return nil
}
