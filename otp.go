package beemafrica

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
)

const otpBaseURL = "https://apiotp.beem.africa"

// Request generates a random OTP and sends it to the provided phone number,application id.
// Requires Mobile number in valid international number format with country code.
// No leading + sign. Example 255712345678. appid is found in beem dashboard.
func (o *Client) Request(number string, appId int) (*http.Response, error) {
	var requestURL = path.Join(otpBaseURL, version, "request")

	body := map[string]any{
		"appId":  appId,
		"msisdn": number,
	}

	bb, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(bb))
	if err != nil {
		return nil, err
	}

	authHeader := generateBasicHeader(o.apiKey, o.apiSecret)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Verify checks to see if the provided OTP matches the pinId provided.
// Returns a Valid 200 OK Response, In Both cases. Look into data for Valid or Invalid OTP
func (o *Client) Verify(pinId string, otp string) (*http.Response, error) {
	var verifyURL = path.Join(otpBaseURL, version, "request")

	body := map[string]any{
		"pinId": pinId,
		"pin":   otp,
	}

	bb, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, verifyURL, bytes.NewBuffer(bb))
	if err != nil {
		return nil, err
	}

	authHeader := generateBasicHeader(o.apiKey, o.apiSecret)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
