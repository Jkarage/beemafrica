package beemafrica

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
)

const (
	airtimeBaseURL = "https://apiairtime.beem.africa"
	topupBaseURL   = "https://apitopup.beem.africa"
)

// Transfer attempts to transfer amount from your account to another account.
// address is the phone number in format 2557135070XX,followed by the amount
// reference is a random number for reference
func (c *Client) Transfer(address string, amount, reference int) (*http.Response, error) {
	var tansferURL = path.Join(airtimeBaseURL, version, "transfer")

	body := map[string]any{
		"dest_addr":    address,
		"amount":       amount,
		"reference_id": reference,
	}

	bb, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, tansferURL, bytes.NewBuffer(bb))
	if err != nil {
		return nil, err
	}

	authHeader := generateBasicHeader(c.apiKey, c.apiSecret)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetBallance retrieves the ballance in your beemafrica account.
func (c *Client) GetBallance() (*http.Response, error) {
	var topupURL = path.Join(airtimeBaseURL, version, "credit-ballance?app_name=AIRTIME")

	authHeader := generateBasicHeader(c.apiKey, c.apiSecret)

	req, err := http.NewRequest(http.MethodGet, topupURL, nil)
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
