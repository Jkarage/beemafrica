package beemafrica

import (
	"context"
	"net/http"
	"path"
)

const (
	airtimeBaseURL = "https://apiairtime.beem.africa"
	topupBaseURL   = "https://apitopup.beem.africa"
)

type AirtimeBallance struct {
	Data struct {
		CreditBalance float64 `json:"credit_bal"`
	} `json:"data"`
}

type AirtimeTransfer struct {
	Code          int    `json:"code"`
	TransactionId int    `json:"transaction_id"`
	Message       string `json:"message"`
}

// Transfer attempts to transfer amount from your account to another account.
// address is the phone number in format 255712345678,followed by the amount
// reference is a random number for reference
func (c *Client) Transfer(ctx context.Context, address string, amount, reference int) (AirtimeTransfer, error) {
	var tansferURL = path.Join(airtimeBaseURL, version, "transfer")
	var resp AirtimeTransfer

	body := map[string]any{
		"dest_addr":    address,
		"amount":       amount,
		"reference_id": reference,
	}

	if err := c.Do(ctx, http.MethodGet, tansferURL, body, resp); err != nil {
		return AirtimeTransfer{}, err
	}

	return resp, nil
}

// GetBallance retrieves the ballance in your beemafrica account.
func (c *Client) GetBallance(ctx context.Context) (AirtimeBallance, error) {
	var topupURL = path.Join(airtimeBaseURL, version, "credit-ballance?app_name=AIRTIME")
	var resp AirtimeBallance

	if err := c.Do(ctx, http.MethodGet, topupURL, nil, resp); err != nil {
		return AirtimeBallance{}, err
	}

	return resp, nil
}
