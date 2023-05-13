# whatsapp-speech-to-text-bot

This [encore.dev](https://encore.dev/) app can receive Twilio webhook events and transcribe audio messages using Open AI's [Whisper](https://github.com/openai/whisper) model, which is deployed separately as a separate [Cloud Run](https://cloud.google.com/run) app.

## Setup

Once deployed, set the twilio webhook to: `<your-app-url>/message/receive`.

By visiting the Sandbox settings in your account:
https://console.twilio.com/us1/develop/sms/try-it-out/whatsapp-learn

## Local testing

Start the app locally:

```
❯ cd whatsapp-speech-to-text-bot

❯ encore run
```

In a separate window, start the Whisper service locally:

```
❯ cd whisper 

❯ make localhost
```

Expose your localhost to the public web [with ngrok](https://www.twilio.com/blog/test-your-webhooks-locally-with-ngrok-html):

```
❯ ngrok http 4000
```

Set your twilio webhook to: `https://<your-forwarding-address>.eu.ngrok.io/message/receive`.

By visiting the Sandbox settings in your account:
https://console.twilio.com/us1/develop/sms/try-it-out/whatsapp-learn

## Limitations

This setup currently deploys the Cloud Run (Whisper) app with [public, unauthenticated access](https://cloud.google.com/run/docs/authenticating/public).
