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

//encore:api public method=POST path=/message/send
func SendMessage(ctx context.Context, p *MessageParams) error {
	client := getTwilioClient()

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(p.To)
	params.SetFrom(p.From)
	params.SetBody(p.Body)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		rlog.Error("error sending message", "err", err)
		return err
	} else {
		response, _ := json.Marshal(*resp)
		rlog.Info("message sent", "response", string(response))
	}
	return nil
}
