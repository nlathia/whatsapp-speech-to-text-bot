package message

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"encore.dev"
	"encore.dev/rlog"
)

// TranscriptionRequest replicates the
// dataclass that is expected by the Python service
type TranscriptionRequest struct {
	MediaUrl string `json:"media_uri"`
}

// TranscriptionResponse replicates the
// dataclass that is returned by the Python service
type TranscriptionResponse struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

// transcriptionServiceUrl returns the URL to the
// Cloud Run function for the current environment
func transcriptionServiceUrl() string {
	if encore.Meta().Environment.Type == encore.EnvProduction {
		return "https://openai-transcribe-dv5eoq6nda-nw.a.run.app"
	} else {
		return "https://openai-transcribe-hbssw3ph2q-nw.a.run.app"
	}
}

// buildRequest returns an http.Request to call the
// Python Cloud Run service
func buildRequest(ctx context.Context, mediaUrl string) (*http.Request, error) {
	if mediaUrl == "" {
		return nil, fmt.Errorf("mediaUrl is nil")
	}

	request, err := json.Marshal(&TranscriptionRequest{
		MediaUrl: mediaUrl,
	})
	if err != nil {
		return nil, err
	}

	url := transcriptionServiceUrl()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// transcribe returns a transcription of a mediaUri by making an http
// request to the openai-transcribe service
func transcribe(ctx context.Context, mediaUri string) (*TranscriptionResponse, error) {
	log := rlog.With("uri", mediaUri)
	req, err := buildRequest(ctx, mediaUri)
	if err != nil {
		log.Error("transcribe failed to build request")
		return nil, err
	}

	// ctx = context.WithTimeout() // @TODO investigate
	client := &http.Client{
		Timeout: time.Second * 60 * 2, // @TODO investigate
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("transcribe failed to Do() request")
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Error("transcribe request failed", "status", resp.Status)
		return nil, fmt.Errorf("transcribe failed with status: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("transcribe body read failed")
		return nil, err
	}

	var result *TranscriptionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Error("transcribe decode failed")
		return nil, err
	}
	return result, nil
}
