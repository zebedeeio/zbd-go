package zebedee

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// BaseURL is https://api.zebedee.io/v0 by default.
type Client struct {
	BaseURL    string
	APIKey     string
	HttpClient *http.Client
}

func New(apikey string) *Client {
	return &Client{
		BaseURL:    "https://api.zebedee.io/v0",
		APIKey:     apikey,
		HttpClient: &http.Client{},
	}
}

type Response struct {
	Success *bool           `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (c *Client) MakeRequest(
	method string,
	path string,
	content interface{},
	response interface{},
) error {
	body := &bytes.Buffer{}
	if content != nil {
		json.NewEncoder(body).Encode(content)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.APIKey)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	var baseResponse Response
	err = json.Unmarshal(responseBody, &baseResponse)
	if err != nil {
		return fmt.Errorf("fail to decode JSON from %s: %s", path, err.Error())
	}

	if resp.StatusCode >= 300 {
		// the API returned a structured error
		if baseResponse.Message != "" {
			return errors.New(baseResponse.Message)
		}

		// an unexpected failure
		return fmt.Errorf("%s returned an error (%d): '%s'",
			path, resp.StatusCode, string(responseBody))
	}

	err = json.Unmarshal(baseResponse.Data, &response)
	if err != nil {
		return fmt.Errorf("Error unmarshaling field \"data\" from API response: %w", err)
	}

	return nil
}

func (c *Client) Wallet() (*Wallet, error) {
	var wallet Wallet
	err := c.MakeRequest("GET", "/wallet", nil, &wallet)
	return &wallet, err
}

// Create Charge: https://api-reference.zebedee.io/#b77ef5ff-477d-4e14-91d0-1713ac06539b
//
// Takes a Charge object containing only
// {ExpiresIn, Amount, Description, InternalID, CallbackURL}
// and overwrites that with the response.
func (c *Client) Charge(params *Charge) (*Charge, error) {
	err := c.MakeRequest("POST", "/charges", params, params)
	return params, err
}

// Get All Charges: https://api-reference.zebedee.io/#cdb9c0d1-76e5-4949-9bb8-e8a52d6aaed3
func (c *Client) ListCharges() ([]Charge, error) {
	var charges []Charge
	err := c.MakeRequest("GET", "/charges", nil, &charges)
	return charges, err
}

// Get Charge Details: https://api-reference.zebedee.io/#a5a2d24c-2a38-44d0-bc00-57598066f1f2
func (c *Client) GetCharge(chargeID string) (*Charge, error) {
	var charge Charge
	err := c.MakeRequest("GET", "/charges/"+chargeID, nil, &charge)
	return &charge, err
}

// Create Withdrawal Request: https://api-reference.zebedee.io/#60cee894-009f-40dc-9cba-e9aec5ce8aa9
//
// Takes an WithdrawalRequest object containing only
// {expiresIn, Amount, Description, InternalID, CallbackURL}
func (c *Client) WithdrawalRequest(params *WithdrawalRequest) (*WithdrawalRequest, error) {
	err := c.MakeRequest("POST", "/withdrawal-requests", params, params)
	return params, err
}

// Get All Withdrawal Requests: https://api-reference.zebedee.io/#bc59c1da-4d5a-49c6-937f-f95d71c940c6
func (c *Client) ListWithdrawalRequests() ([]WithdrawalRequest, error) {
	var wr []WithdrawalRequest
	err := c.MakeRequest("GET", "/withdrawal-requests", nil, &wr)
	return wr, err
}

// Get Withdrawal Request Details: https://api-reference.zebedee.io/#12aea552-0b8d-4562-a84b-a890d4f17a32
func (c *Client) GetWithdrawalRequest(wrequestID string) (*WithdrawalRequest, error) {
	var wr WithdrawalRequest
	err := c.MakeRequest("GET", "/withdrawal-requests/"+wrequestID, nil, &wr)
	return &wr, err
}

// Pay Invoice: https://api-reference.zebedee.io/#04dace34-06f5-4c2f-9215-5870205098d5
//
// Takes a Payment object containing only {Description, InternalID, Invoice}
// and overwrites that with the response.
func (c *Client) Pay(params *Payment) (*Payment, error) {
	err := c.MakeRequest("POST", "/payments", params, params)
	return params, err
}

// Get All Payments: https://api-reference.zebedee.io/#08ea69cc-dd6f-4381-a489-18004b911f96
func (c *Client) ListPayments() ([]Payment, error) {
	var payments []Payment
	err := c.MakeRequest("GET", "/payments", nil, &payments)
	return payments, err
}

// Get Payment Details: https://api-reference.zebedee.io/#244ebe9f-6c4d-4162-a805-9a0e8955b20d
func (c *Client) GetPayment(paymentID string) (*Payment, error) {
	var payment Payment
	err := c.MakeRequest("GET", "/payments/"+paymentID, nil, &payment)
	return &payment, err
}

// Get API Production IPs: https://api-reference.zebedee.io/#c7e18276-6935-4cca-89ae-ad949efe9a6a
func (c *Client) GetProductionIPs() ([]string, error) {
	var ips struct {
		IPs []string `json:"ips"`
	}
	err := c.MakeRequest("GET", "/prod-ipds", nil, &ips)
	return ips.IPs, err
}
