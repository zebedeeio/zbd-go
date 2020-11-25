package zebedee

import (
	"bytes"
	"encoding/json"
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
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%s%s returned an error: '%s'",
			c.BaseURL, path, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(response)
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
func (c *Client) ListCharges(params *Charge) ([]Charge, error) {
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

// Pay Invoice: https://api-reference.zebedee.io/#04dace34-06f5-4c2f-9215-5870205098d5
//
// Takes a Payment object containing only {Description, internalID, Invoice}
// and overwrites that with the response.
func (c *Client) Pay(params *Payment) (*Payment, error) {
	err := c.MakeRequest("POST", "/", params, params)
	return params, err
}

// Get All Payments: https://api-reference.zebedee.io/#08ea69cc-dd6f-4381-a489-18004b911f96
func (c *Client) ListPayments(params *Payment) ([]Payment, error) {
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
