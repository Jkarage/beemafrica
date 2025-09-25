package beemafrica

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"
)

var ErrInvalidCreds = errors.New("failed to load accounts apikey or secretkey")

const smsBaseURL = "https://apisms.beem.africa"

type Recipients struct {
	RecipientId int    `json:"recipient_id"`
	Destination string `json:"dest_addr"`
}

type SMSInput struct {
	SenderID     string       `json:"source_addr"`
	ScheduleTime string       `json:"schedule_time"`
	Encoding     string       `json:"encoding"`
	Message      string       `json:"message"`
	Recipients   []Recipients `json:"recipients"`
}

type SMSResponse struct {
	Successful bool   `jso:"successful"`
	RequestId  int    `json:"request_id"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Valid      int    `json:"valid"`
	InValid    int    `json:"invalid"`
	Duplicates int    `json:"duplicates"`
}

type SenderName struct {
	ID            string `json:"id"`
	SenderID      string `json:"senderid"`
	SampleContent string `json:"sample_content"`
	Status        string `json:"status"`
	Created       string `json:"created"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	TotalItems  int   `json:"totalItems"`
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	TotalPages  int   `json:"totalPages"`
	StartPage   int   `json:"startPage"`
	EndPage     int   `json:"endPage"`
	StartIndex  int   `json:"startIndex"`
	EndIndex    int   `json:"endIndex"`
	Pages       []int `json:"pages"`
}

type SenderNames struct {
	Data       []SenderName   `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

type SMSBalance struct {
	Data struct {
		CreditBalance float64 `json:"credit_balance"`
	} `json:"data"`
}

type SenderNameInput struct {
	Data struct {
		SenderID      string `json:"senderid"`
		SampleContent string `json:"sample_content"`
		ID            string `json:"id"`
		Status        string `json:"status"`
	} `json:"data"`
}

// SendSMS sends request to beemafrica to send a message, with given details.
// the message, a slice of recipients, and a scheduled time value.
// time format is  GMT+0 timezone,(yyyy-mm-dd hh:mm).
// send now scheduled_time is ""
func (c *Client) SendSMS(ctx context.Context, message string, recipients []string, schedule_time, senderID string) (*SMSResponse, error) {
	var resp = &SMSResponse{}

	sendSMSURL, err := url.JoinPath(smsBaseURL, version, "send")
	if err != nil {
		return nil, err
	}

	body := SMSInput{
		SenderID:     senderID,
		ScheduleTime: schedule_time,
		Encoding:     "0",
		Message:      message,
		Recipients:   make([]Recipients, 0),
	}

	for i, r := range recipients {
		body.Recipients = append(body.Recipients, Recipients{
			RecipientId: i + 1,
			Destination: r,
		})
	}

	if err := c.Do(ctx, http.MethodPost, sendSMSURL, &body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetBallance request for the sms ballance for a particular account
// If the error is nil, the response of type *SMSBalance will be returned
func (c *Client) GetSMSBallance(ctx context.Context) (*SMSBalance, error) {
	var resp = &SMSBalance{}

	ballanceURL, err := url.JoinPath(smsBaseURL, "public", version, "vendors", "balance")
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, http.MethodGet, ballanceURL, nil, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// RequestSenderID queues a request to beem for a specific senderid.
// Response will be obtained, later through mail.
func (c *Client) RequestSenderID(ctx context.Context, id, idContent string) (*SenderName, error) {
	var resp = &SenderName{}
	var senderURL = path.Join(smsBaseURL, "public", version, "sender-names")

	body := map[string]string{
		"senderid":       id,
		"sample_content": idContent,
	}

	if err := c.Do(ctx, http.MethodPost, senderURL, body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetSenderNames retrieves sendernames available in your account.
func (c *Client) GetSenderNames(ctx context.Context) (*SenderNames, error) {
	var senderURL = path.Join(smsBaseURL, "public", version, "sender-names")
	var resp = &SenderNames{}

	if err := c.Do(ctx, http.MethodGet, senderURL, nil, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
