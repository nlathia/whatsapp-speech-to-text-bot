package message

import (
	"strconv"

	"github.com/twilio/twilio-go"
)

type TwilioMessage struct {
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

var secrets struct {
	TwilioAccountSid string
	TwilioAuthToken  string
}

func getTwilioClient() *twilio.RestClient {
	return twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: secrets.TwilioAccountSid,
		Password: secrets.TwilioAuthToken,
	})
}

// parseBool converts a string to a bool, if the key to that
// string exists in values (otherwise, defaults to returning false)
func parseBool(key string, values map[string]string) (bool, error) {
	value, exists := values[key]
	if !exists {
		return false, nil
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return result, nil
}

// parseInt converts a string to a int, if the key to that
// string exists in values (otherwise, defaults to returning 0)
func parseInt(key string, values map[string]string) (int, error) {
	value, exists := values[key]
	if !exists {
		return 0, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// firstValues transforms a map[string][]string into a
// map[string]string by taking the first value from each key
// in the input map
func firstValues(form map[string][]string) map[string]string {
	result := map[string]string{}
	for key, values := range form {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

// formToMessage returns a TwilioMessage that is created by
// parsing the fields in the webhook's form values
func formToMessage(form map[string][]string) (*TwilioMessage, error) {
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

	return &TwilioMessage{
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
