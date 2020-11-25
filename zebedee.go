package zebedee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func (c *Client) Charge(params *Charge) (*Charge, error) {
	err := c.MakeRequest("POST", "/charges", params, params)
	return params, err
}

func (c *Client) ListCharges(params *Charge) ([]Charge, error) {
	var charges []Charge
	err := c.MakeRequest("GET", "/charges", nil, &charges)
	return charges, err
}

func (c *Client) GetCharge(chargeID string) (*Charge, error) {
	var charge Charge
	err := c.MakeRequest("GET", "/charges/"+chargeID, nil, &charge)
	return &charge, err
}

func (c *Client) Pay(params *Payment) (*Payment, error) {
	err := c.MakeRequest("POST", "/", params, params)
	return params, err
}

func (c *Client) ListPayments(params *Payment) ([]Payment, error) {
	var payments []Payment
	err := c.MakeRequest("GET", "/payments", nil, &payments)
	return payments, err
}

func (c *Client) GetPayment(paymentID string) (*Payment, error) {
	var payment Payment
	err := c.MakeRequest("GET", "/payments/"+paymentID, nil, &payment)
	return &payment, err
}
