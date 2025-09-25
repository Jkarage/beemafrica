package beemafrica

import (
	"context"
	"net/http"
	"path"
)

const otpBaseURL = "https://apiotp.beem.africa"

// OTPResponse represents the common structure for OTP API responses
type OTP struct {
	Data Data `json:"data"`
}
type Data struct {
	Message *Message `json:"message,omitempty"`
	PinID   string   `json:"pinId,omitempty"`
	Status  string   `json:"status,omitempty"`
}

// ResponseMessage represents the message object with code and description
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Request generates a random OTP and sends it to the provided phone number.
// Requires mobile number in valid international format (no leading +), e.g., "255712345678"
func (c *Client) Request(ctx context.Context, number string, appId int) (*OTP, error) {
	var requestURL = path.Join(otpBaseURL, version, "request")
	var resp = &OTP{}

	body := map[string]any{
		"appId":  appId,
		"msisdn": number,
	}

	if err := c.Do(ctx, http.MethodPost, requestURL, body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Verify checks to see if the provided OTP matches the pinId provided.
// Returns a Valid 200 OK Response, In Both cases. Look into data for Valid or Invalid OTP
func (c *Client) Verify(ctx context.Context, pinId string, otp string) (*OTP, error) {
	var verifyURL = path.Join(otpBaseURL, version, "request")
	var resp = &OTP{}

	body := map[string]any{
		"pinId": pinId,
		"pin":   otp,
	}

	if err := c.Do(ctx, http.MethodPost, verifyURL, body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
