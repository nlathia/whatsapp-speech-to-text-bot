package whisper

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

const (
	// cloudRunURL is the Cloud Run service that is running Python
	// (implemented in this directory)
	cloudRunURL   = "https://openai-transcribe-%v-nw.a.run.app"
	cloudRunLocal = "http://localhost:9090/"
)

var secrets struct {
	CloudRunID string
}

// WhisperRequest replicates the
// dataclass that is expected by the Python service
type WhisperRequest struct {
	MediaUrl string `json:"media_uri"`
}

// WhisperResponse replicates the
// dataclass that is returned by the Python service
type WhisperResponse struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

func isLocal() bool {
	if encore.Meta().Environment.Type == encore.EnvDevelopment {
		return encore.Meta().Environment.Cloud == encore.CloudLocal
	}
	return false
}

func getURL() string {
	if isLocal() {
		return cloudRunLocal
	}
	return fmt.Sprintf(cloudRunURL, secrets.CloudRunID)
}

// buildRequest returns an http.Request to call the
// Python Cloud Run service
func (w *WhisperRequest) buildRequest(ctx context.Context) (*http.Request, error) {
	url := getURL()
	request, err := json.Marshal(&WhisperRequest{
		MediaUrl: w.MediaUrl,
	})
	if err != nil {
		return nil, err
	}

	rlog.Info("building request", "service", url)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Transcribe returns a transcription of a mediaUri by making an http
// request to the openai-transcribe service
//
//encore:api private method=POST path=/whisper
func Transcribe(ctx context.Context, params *WhisperRequest) (*WhisperResponse, error) {
	req, err := params.buildRequest(ctx)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 60 * 2,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("transcribe failed with status: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result *WhisperResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
