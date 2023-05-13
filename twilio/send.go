package twilio

import (
	"context"
	"encoding/json"
	"errors"

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

func (m *MessageParams) Validate() error {
	switch {
	case m.From == "":
		return errors.New("from field is required")
	case m.To == "":
		return errors.New("to field is required")
	case m.Body == "":
		return errors.New("body field is required")
	}
	return nil
}

// SendMessage sends a message via Twilio
//
//encore:api private method=POST path=/message/send
func SendMessage(ctx context.Context, p *MessageParams) error {
	if err := p.Validate(); err != nil {
		return err
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: secrets.TwilioAccountSid,
		Password: secrets.TwilioAuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(p.To)
	params.SetFrom(p.From)
	params.SetBody(p.Body)
	resp, err := client.Api.CreateMessage(params)
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
