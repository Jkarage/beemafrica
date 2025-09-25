package beemafrica

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"path"
)

var ErrInvalidCreds = errors.New("failed to load accounts apikey or secretkey")

const smsBaseURL = "https://apisms.beem.africa"

// SendSMS sends request to beemafrica to send a message, with given details.
// the message, a slice of recipients, and a scheduled time value.
// time format is  GMT+0 timezone,(yyyy-mm-dd hh:mm).
// send now scheduled_time is ""
func (c *Client) SendSMS(message string, recipients []string, schedule_time, senderID string) (*http.Response, error) {
	var (
		resp       *http.Response
		sendSMSURL = path.Join(smsBaseURL, version, "send")
	)

	for i, r := range recipients {
		body := map[string]interface{}{
			"source_addr":   senderID,
			"schedule_time": schedule_time,
			"encoding":      "0",
			"message":       message,
			"recipients": []map[string]any{
				{
					"recipient_id": i + 1,
					"dest_addr":    r,
				},
			},
		}

		// convert the body to json
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		// Create a new request
		req, err := http.NewRequest(http.MethodPost, sendSMSURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, err
		}

		authHeader := generateBasicHeader(c.apiKey, c.apiSecret)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

// GetBallance request for the sms ballance for a particular account
// If the error is nil, the response of type *http.Response will be returned
func (c *Client) GetSMSBallance() (*http.Response, error) {
	var (
		resp        *http.Response
		ballanceURL = path.Join(smsBaseURL, "public", version, "vendors", "ballance")
	)

	// Create a new request
	req, err := http.NewRequest(http.MethodGet, ballanceURL, nil)
	if err != nil {
		return resp, err
	}

	authHeader := generateBasicHeader(c.apiKey, c.apiSecret)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RequestSenderID queues a request to beem for a specific senderid.
// Response will be obtained, later through mail.
func (c *Client) RequestSenderID(id, idContent string) (*http.Response, error) {
	var senderURL = path.Join(smsBaseURL, "public", version, "sender-names")

	body := map[string]string{
		"senderid":       id,
		"sample_content": idContent,
	}

	mb, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, senderURL, bytes.NewBuffer(mb))
	if err != nil {
		return nil, err
	}

	authHeader := generateBasicHeader(c.apiKey, c.apiSecret)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetSenderNames retrieves sendernames available in your account.
func (c *Client) GetSenderNames() (*http.Response, error) {
	var senderURL = path.Join(smsBaseURL, "public", version, "sender-names")

	authHeader := generateBasicHeader(c.apiKey, c.apiSecret)

	req, err := http.NewRequest(http.MethodGet, senderURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
