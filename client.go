package beemafrica

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func NewClient(log Logger, apiKey string, apiSecret string, options ...func(cln *Client)) *Client {
	cln := Client{
		log:       log,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		http:      &defaultClient,
	}

	for _, option := range options {
		option(&cln)
	}

	if cln.apiKey == "" || cln.apiSecret == "" {
		fmt.Println(ErrInvalidCreds)
		os.Exit(-1)
	}

	return &cln
}

func WithClient(http *http.Client) func(cln *Client) {
	return func(cln *Client) {
		cln.http = http
	}
}

func (cln *Client) Do(ctx context.Context, method string, endpoint string, body any, v any) error {
	resp, err := do(ctx, cln, method, endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("client: copy error: %w", err)
	}

	switch d := v.(type) {
	case *string:
		*d = string(data)

	default:
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("client: response: %s, decoding error: %w ", string(data), err)
		}
	}

	return nil
}

func do(ctx context.Context, cln *Client, method string, endpoint string, body any) (*http.Response, error) {
	var statusCode int

	cln.log(ctx, "do: rawRequest: started", "method", method, "endpoint", endpoint)
	defer func() {
		cln.log(ctx, "do: rawRequest: completed", "status", statusCode)
	}()

	var b bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&b).Encode(body); err != nil {
			return nil, fmt.Errorf("encoding: error: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, &b)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	token := generateBasicHeader(cln.apiKey, cln.apiSecret)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))

	resp, err := cln.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: error: %w", err)
	}

	// Assign for logging the status code at the end of the function call.
	statusCode = resp.StatusCode

	switch statusCode {
	case http.StatusOK, http.StatusNoContent:
		return resp, nil

	default:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("readall: error: %w", err)
		}

		if err := json.Unmarshal(data, &err); err != nil {
			return nil, fmt.Errorf("decoding: response: %s, error: %w ", string(data), err)
		}

		return nil, fmt.Errorf("error: response: %s", err)
	}
}
