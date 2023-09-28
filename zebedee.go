package zebedee

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ZBDOauth struct {
	ClientID    string
	Secret      string
	RedirectURI string
	State       string
	Scope       string
}

// BaseURL is https://api.zebedee.io/v0 by default.
type Client struct {
	BaseURL    string
	APIKey     string
	HttpClient *http.Client
	Oauth      *ZBDOauth
}

func newFetchTokenBody(c *Client, code string, codeVerifier string) FetchTokenBody {
	return FetchTokenBody{
		ClientID:     c.Oauth.ClientID,
		ClientSecret: c.Oauth.Secret,
		Code:         code,
		CodeVerifier: codeVerifier,
		GrantType:    "authorization_code",
		RedirectURI:  c.Oauth.RedirectURI,
	}
}

func newFetchRefresh(c *Client, refreshToken string) FetchRefresh {
	return FetchRefresh{
		ClientID:     c.Oauth.ClientID,
		ClientSecret: c.Oauth.Secret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	}
}

func New(apikey string, oauth *ZBDOauth) *Client {
	return &Client{
		BaseURL:    "https://api.zebedee.io/v0",
		APIKey:     apikey,
		HttpClient: &http.Client{},
		Oauth:      oauth,
	}
}

func NewOauth(client_id string, secret string, redirect_uri string, state string, scope string) *ZBDOauth {
	return &ZBDOauth{
		ClientID:    client_id,
		Secret:      secret,
		RedirectURI: redirect_uri,
		State:       state,
		Scope:       scope,
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
	headers http.Header,
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
	if headers != nil {
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

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
	err := c.MakeRequest("GET", "/wallet", nil, &wallet, nil)
	return &wallet, err
}

// Create Charge: https://api-reference.zebedee.io/#b77ef5ff-477d-4e14-91d0-1713ac06539b
//
// Takes a Charge object containing only
// {ExpiresIn, Amount, Description, InternalID, CallbackURL}
// and overwrites that with the response.
func (c *Client) Charge(params *Charge) (*Charge, error) {
	err := c.MakeRequest("POST", "/charges", params, params, nil)
	return params, err
}

// DecodeCharge decodes charge information using the provided invoice string.
// It makes a POST request to the "/charges" API endpoint with the given parameters
// and returns the decoded charge details.
//
// Parameters:
//   - param: A pointer to a DecodeChargeOptionsType struct containing the invoice string.
//
// Returns:
//   - A pointer to a DecodeChargeResponseType struct representing the decoded charge response.
//   - An error if the request fails or an error occurs during response processing.
//
// Example:
//
//	client := NewClient(apiKey)
//	invoice := "lnbc123456789" // Replace with the actual invoice
//	options := &DecodeChargeOptionsType{Invoice: invoice}
//	response, err := client.DecodeCharge(options)
//	if err != nil {
//	  fmt.Println("Error:", err)
//	  return
//	}
//	fmt.Println("Decoded Charge:", response.Data)
func (c *Client) DecodeCharge(param *DecodeChargeOptionsType) (*DecodeChargeResponseType, error) {
	var res DecodeChargeResponseType
	err := c.MakeRequest("POST", "/charges", param, &res, nil)
	return &res, err
}

// Get All Charges: https://api-reference.zebedee.io/#cdb9c0d1-76e5-4949-9bb8-e8a52d6aaed3
func (c *Client) ListCharges() ([]Charge, error) {
	var charges []Charge
	err := c.MakeRequest("GET", "/charges", nil, &charges, nil)
	return charges, err
}

// Get Charge Details: https://api-reference.zebedee.io/#a5a2d24c-2a38-44d0-bc00-57598066f1f2
func (c *Client) GetCharge(chargeID string) (*Charge, error) {
	var charge Charge
	err := c.MakeRequest("GET", "/charges/"+chargeID, nil, &charge, nil)
	return &charge, err
}

// Static charges endpoints

// CreateStaticCharge creates a new static charge with the provided parameters.
// It makes a POST request to the "/static-charges" API endpoint using the given parameters,
// and returns the response containing the newly created static charge details.
//
// Parameters:
//   - param: An instance of the StaticChargeOptionsType struct containing charge creation options.
//
// Returns:
//   - A pointer to a StaticChargeDataResponseType struct representing the response
//     containing the details of the newly created static charge.
//   - An error if the request fails or an error occurs during response processing.
//
// Example:
//
//	client := NewClient(apiKey)
//	chargeOptions := StaticChargeOptionsType{
//	  MinAmount:      "1000", // Replace with desired values
//	  MaxAmount:      "50000",
//	  Description:    "Sample charge",
//	  InternalID:     "charge123",
//	  CallbackURL:    "https://example.com/callback",
//	  SuccessMessage: "Charge successful",
//	}
//	response, err := client.CreateStaticCharge(chargeOptions)
//	if err != nil {
//	  fmt.Println("Error:", err)
//	  return
//	}
//	fmt.Println("Created Charge ID:", response.Data.ID)
func (c *Client) CreateStaticCharge(param StaticChargeOptionsType) (*StaticChargeDataResponseType, error) {
	var res StaticChargeDataResponseType
	err := c.MakeRequest("POST", "/static-charges", param, &res, nil)
	return &res, err
}

// GetStaticCharge retrieves the details of a static charge by its ID.
//
// The function makes a GET request to the API's "static-charges" endpoint using the provided
// staticChargeID as a parameter to fetch the information about the static charge.
//
// Parameters:
//   - staticChargeID: The unique identifier of the static charge to retrieve.
//
// Returns:
//   - *StaticChargeDataResponseType: A pointer to the response containing the static charge details.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	chargeID := "sample-charge-id"
//	response, err := client.GetStaticCharge(chargeID)
//	if err != nil {
//	  fmt.Println("Error fetching static charge:", err)
//	  return
//	}
//	fmt.Println("Static charge data:", response.Data)
func (c *Client) GetStaticCharge(staticChargeID string) (*StaticChargeDataResponseType, error) {
	var res StaticChargeDataResponseType
	err := c.MakeRequest("GET", "/static-charges/"+staticChargeID, nil, &res, nil)
	return &res, err
}

// UpdateStaticCharge updates the properties of a static charge identified by its ID.
//
// The function sends a PATCH request to the API's "static-charges" endpoint with the provided
// staticChargeID as a parameter to update the properties of the specified static charge.
// The updated properties are specified in the param parameter, which should be of type StaticChargeOptionsType.
//
// Parameters:
//   - staticChargeID: The unique identifier of the static charge to update.
//   - param: A StaticChargeOptionsType object containing the updated properties of the static charge.
//
// Returns:
//   - *StaticChargeDataResponseType: A pointer to the response containing the updated static charge details.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	chargeID := "sample-charge-id"
//	updatedCharge, err := client.UpdateStaticCharge(chargeID, StaticChargeOptionsType{
//	  MinAmount:      "2000",
//	  MaxAmount:      "60000",
//	  Description:    "Updated charge",
//	  InternalID:     "charge456",
//	  CallbackURL:    "https://example.com/update-callback",
//	  SuccessMessage: "Updated charge successful",
//	})
//	if err != nil {
//	  fmt.Println("Error updating static charge:", err)
//	  return
//	}
//	fmt.Println("Updated static charge data:", updatedCharge.Data)
func (c *Client) UpdateStaticCharge(staticChargeID string, param StaticChargeOptionsType) (*StaticChargeDataResponseType, error) {
	var res StaticChargeDataResponseType
	err := c.MakeRequest("PATCH", "/static-charges/"+staticChargeID, param, &res, nil)
	return &res, err
}

// Lightening Address

// SendLightningAddressPayment initiates a lightning payment to the specified lightning address.
//
// The function sends a POST request to the API's "ln-address/send-payment/" endpoint with the provided
// payment options specified in the param parameter to send a lightning payment to the specified lightning address.
//
// Parameters:
//   - param: A SendLightningAddressPaymentOptionsType object containing the payment options and details.
//
// Returns:
//   - *SendLightningAddressPaymentDataResponseType: A pointer to the response containing the payment details.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	paymentOptions := SendLightningAddressPaymentOptionsType{
//	  LnAddress:   "lnbc123...",
//	  Amount:      "1000",
//	  Comment:     "Payment for goods",
//	  CallbackUrl: "https://example.com/payment-callback",
//	  InternalID:  "payment123",
//	}
//	paymentResponse, err := client.SendLightningAddressPayment(paymentOptions)
//	if err != nil {
//	  fmt.Println("Error sending lightning address payment:", err)
//	  return
//	}
//	fmt.Println("Payment ID:", paymentResponse.Data.ID)
func (c *Client) SendLightningAddressPayment(param SendLightningAddressPaymentOptionsType) (*SendLightningAddressPaymentDataResponseType, error) {
	var res SendLightningAddressPaymentDataResponseType
	err := c.MakeRequest("POST", "/ln-address/send-payment/", param, &res, nil)
	return &res, err
}

// ValidateLightningAddress validates a Lightning Network (LN) address.
//
// The function makes a GET request to the API's "ln-address/validate" endpoint,
// appending the provided lightningAddress to the URL, to check the validity of
// the given LN address. The response contains information about whether the address
// is valid and additional metadata associated with the address.
//
// Parameters:
//   - lightningAddress: The Lightning Network address to validate.
//
// Returns:
//   - *ValidateLightningAddressDataResponseType: A pointer to the response containing
//     the validation result and associated metadata.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	lnAddress := "lnbc1..."
//	response, err := client.ValidateLightningAddress(lnAddress)
//	if err != nil {
//	  fmt.Println("Error validating LN address:", err)
//	  return
//	}
//	fmt.Println("LN address validation result:", response.Data.Valid)
func (c *Client) ValidateLightningAddress(lightningAddress string) (*ValidateLightningAddressDataResponseType, error) {
	var res ValidateLightningAddressDataResponseType
	err := c.MakeRequest("GET", "/ln-address/validate/"+lightningAddress, nil, &res, nil)
	return &res, err
}

// CreateChargeForLightningAddress creates a charge for a Lightning Network address.
//
// This function makes a POST request to the API's "ln-address/fetch-charge" endpoint
// using the provided CreateChargeFromLightningAddressOptionsType parameters to create a
// charge for the specified Lightning Network address.
//
// Parameters:
//   - params: CreateChargeFromLightningAddressOptionsType containing the required parameters
//     for creating the charge.
//
// Returns:
//   - *FetchChargeFromLightningAddressDataResponseType: A pointer to the response containing
//     the charge details for the Lightning Network address.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	chargeParams := CreateChargeFromLightningAddressOptionsType{
//	  Amount:      "10000",
//	  LNAddress:   "lnaddress123",
//	  Description: "Charge for LN address",
//	}
//	response, err := client.CreateChargeForLightningAddress(chargeParams)
//	if err != nil {
//	  fmt.Println("Error creating charge for LN address:", err)
//	  return
//	}
//	fmt.Println("Charge details:", response.Data)
func (c *Client) CreateChargeForLightningAddress(params CreateChargeFromLightningAddressOptionsType) (*FetchChargeFromLightningAddressDataResponseType, error) {
	var res FetchChargeFromLightningAddressDataResponseType
	err := c.MakeRequest("POST", "/ln-address/fetch-charge", params, &res, nil)
	return &res, err
}

// Keysend endpoints:

// SendKeysendPayment initiates a keysend payment to a Lightning Network node.
//
// This function makes a POST request to the API's "/keysend-payment" endpoint using the
// provided KeysendOptionsType parameters to initiate a keysend payment to the specified
// Lightning Network pubkey.
//
// Parameters:
//   - params: KeysendOptionsType containing the required parameters for initiating the keysend payment.
//
// Returns:
//   - *KeysendDataResponseType: A pointer to the response containing information about the initiated keysend payment.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	keysendParams := KeysendOptionsType{
//	  Amount:      "1000",
//	  Pubkey:      "02abcd...",
//	  TLVRecords:  "my-tlv-records",
//	  Metadata:    "additional-metadata",
//	  CallbackURL: "https://example.com/callback",
//	}
//	response, err := client.SendKeysendPayment(keysendParams)
//	if err != nil {
//	  fmt.Println("Error sending keysend payment:", err)
//	  return
//	}
//	fmt.Println("Keysend payment details:", response.Data)
func (c *Client) SendKeysendPayment(params KeysendOptionsType) (*KeysendDataResponseType, error) {
	var res KeysendDataResponseType
	err := c.MakeRequest("POST", "/keysend-payment", params, &res, nil)
	return &res, err
}

// Utils functions

// GetBtcUsdExchangeRate retrieves the current BTC to USD exchange rate.
//
// This function makes a GET request to the API's "/btcusd" endpoint to fetch the
// current BTC to USD exchange rate and its associated timestamp.
//
// Returns:
//   - *BTCUSDDataResponseType: A pointer to the response containing the BTC to USD exchange rate
//     and timestamp.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	response, err := client.GetBtcUsdExchangeRate()
//	if err != nil {
//	  fmt.Println("Error fetching BTC to USD exchange rate:", err)
//	  return
//	}
//	fmt.Printf("BTC to USD exchange rate: %s\n", response.Data.BTCUSDPrice)
//	fmt.Printf("Exchange rate timestamp: %s\n", response.Data.BTCUSDTimestamp)
func (c *Client) GetBTCUSDExchangeRate() (*BTCUSDDataResponseType, error) {
	var res BTCUSDDataResponseType
	err := c.MakeRequest("GET", "/btcusd", nil, &res, nil)
	return &res, err
}

// IsSupportedRegion checks if the specified IP address is from a supported region.
//
// This function makes a GET request to the API's "/is-supported-region" endpoint
// using the provided IP address as a parameter to determine if the region is supported.
//
// Parameters:
//   - ipAddress: The IP address to check for support.
//
// Returns:
//   - *SupportedRegionDataResponseType: A pointer to the response indicating whether
//     the IP address is from a supported region.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	ipAddress := "127.0.0.1"
//	response, err := client.IsSupportedRegion(ipAddress)
//	if err != nil {
//	  fmt.Println("Error checking supported region:", err)
//	  return
//	}
//	fmt.Println("Is supported region:", response.Data.IsSupported)
func (c *Client) IsSupportedRegion(ipAddress string) (*SupportedRegionDataResponseType, error) {
	var res SupportedRegionDataResponseType
	err := c.MakeRequest("GET", "/is-supported-region/"+ipAddress, nil, &res, nil)
	return &res, err
}

// GetZBDProdIps retrieves the list of ZBD production IP addresses.
//
// This function makes a GET request to the API's "/prod-ips" endpoint to fetch
// the list of ZBD production IP addresses.
//
// Returns:
//   - *ProdIPSDataResponseType: A pointer to the response containing the list of ZBD production IP addresses.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	response, err := client.GetZBDProdIps()
//	if err != nil {
//	  fmt.Println("Error fetching ZBD production IPs:", err)
//	  return
//	}
//	fmt.Println("ZBD production IPs:", response.Data.IPS)
func (c *Client) GetZBDProdIps() (*ProdIPSDataResponseType, error) {
	var res ProdIPSDataResponseType
	err := c.MakeRequest("GET", "/prod-ips", nil, &res, nil)
	return &res, err
}

// Internal transfer functions

// InternalTransfer initiates an internal transfer of funds between two user wallets.
//
// This function makes a POST request to the API's "/internal-transfer" endpoint using the provided
// InternalTransferOptionsType parameters to initiate the transfer of funds from the sender wallet
// to the receiver wallet.
//
// Parameters:
//   - param: InternalTransferOptionsType containing the required parameters for the internal transfer.
//     This includes the amount to transfer and the receiver wallet ID.
//
// Returns:
//   - *InternalTransferDataResponseType: A pointer to the response containing details about the internal transfer.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	transferParams := InternalTransferOptionsType{
//	  Amount:           "5000",
//	  ReceiverWalletId: "receiver-wallet-id",
//	}
//	response, err := client.InternalTransfer(transferParams)
//	if err != nil {
//	  fmt.Println("Error initiating internal transfer:", err)
//	  return
//	}
//	fmt.Println("Internal transfer details:", response.Data)
func (c *Client) InternalTransfer(param InternalTransferOptionsType) (*InternalTransferDataResponseType, error) {
	var res InternalTransferDataResponseType
	err := c.MakeRequest("POST", "/internal-transfer", param, &res, nil)
	return &res, err
}

// OAuth functions

// CreateAuthUrl generates an authorization URL for OAuth2 authentication.
//
// This function constructs an authorization URL using the provided client configuration
// and generates a state, code verifier, and code challenge for use in the OAuth2 flow.
//
// Returns:
//   - string: The authorization URL that users can visit to initiate the OAuth2 flow.
//   - error: An error if there's an issue generating the authorization URL.
//
// Example usage:
//
//	oauth := NewOauth("your_client_id", "your_secret", "your_redirect_uri", "your_state", "your_scope")
//	client := New("your_api_key", oauth)
//	authURL, err := client.CreateAuthUrl()
//	if err != nil {
//	  fmt.Println("Error generating auth URL:", err)
//	  return
//	}
//	fmt.Println("Authorization URL:", authURL)
func (c *Client) CreateAuthUrl() (string, error) {
	state := RandomString(30)
	codeVerifier := RandomString(43)
	codeChallenge, _ := S256Challenge(codeVerifier)
	codeChallengeMethod := "S256"

	baseURL, _ := url.Parse(c.BaseURL)
	baseURL.Path = "/v1/oauth2/authorize"
	values := url.Values{}
	values.Add("client_id", c.Oauth.ClientID)
	values.Add("response_type", "code")
	values.Add("redirect_uri", c.Oauth.RedirectURI)
	values.Add("code_challenge_method", codeChallengeMethod)
	values.Add("code_challenge", codeChallenge)
	values.Add("scope", c.Oauth.Scope)
	values.Add("state", state)

	baseURL.RawQuery = values.Encode()
	authURL := baseURL.String()

	return authURL, nil
}

// FetchToken retrieves an access token from the OAuth2 authorization code.
//
// This function makes a POST request to the API's "/v1/oauth2/token" endpoint using the provided
// authorization code to exchange it for an access token.
//
// Parameters:
//   - code: The authorization code obtained from the OAuth2 authorization process.
//
// Returns:
//   - FetchAccessTokenRes: The response containing the fetched access token and related information.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	code := "authorization-code-obtained-from-authorization-process"
//	tokenResponse, err := client.FetchToken(code)
//	if err != nil {
//	  fmt.Println("Error fetching access token:", err)
//	  return
//	}
//	fmt.Println("Access Token:", tokenResponse.AccessToken)
//	fmt.Println("Expires In:", tokenResponse.ExpiresIn)
//	// Handle other token response data
func (c *Client) FetchToken(code string) (FetchAccessTokenRes, error) {
	var res FetchAccessTokenRes
	codeVerifier := RandomString(43)
	_, verifier := S256Challenge(codeVerifier)
	body := newFetchTokenBody(c, code, verifier)
	err := c.MakeRequest("POST", "/v1/oauth2/token", body, &res, nil)
	return res, err
}

// RefreshToken refreshes an access token using a refresh token.
//
// This function makes a POST request to the API's "/v1/oauth2/token" endpoint
// using the provided refresh token to request a new access token.
//
// Parameters:
//   - refreshToken: The refresh token used to obtain a new access token.
//
// Returns:
//   - FetchPostRes: The response containing the refreshed access token details.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	refreshToken := "your-refresh-token"
//	response, err := client.RefreshToken(refreshToken)
//	if err != nil {
//	  fmt.Println("Error refreshing access token:", err)
//	  return
//	}
//	fmt.Println("Refreshed access token:", response.AccessToken)
func (c *Client) RefreshToken(refreshToken string) (FetchPostRes, error) {
	var res FetchPostRes
	body := newFetchRefresh(c, refreshToken)
	err := c.MakeRequest("POST", "/v1/oauth2/token", body, &res, nil)
	return res, err
}

// GetUserData retrieves user data using the provided user token.
//
// This function makes a GET request to the API's "/v1/oauth2/user" endpoint
// to fetch user data for the authenticated user.
//
// Parameters:
//   - token: A string containing the user token.
//
// Returns:
//   - ZBDUserData: The user data fetched from the API.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	userData, err := client.GetUserData("user-token-here")
//	if err != nil {
//	  fmt.Println("Error fetching user data:", err)
//	  return
//	}
//	fmt.Println("User data:", userData)
func (c *Client) GetUserData(token string) (ZBDUserData, error) {
	var res ZBDUserData
	headers := make(http.Header)
	headers.Set("usertoken", token)

	err := c.MakeRequest("GET", "/v1/oauth2/user", nil, &res, headers)

	return res, err
}

// GetUserWalletData retrieves the wallet data for the authenticated user.
//
// This function makes a GET request to the API's "/v1/oauth2/wallet" endpoint using the provided
// user token in the headers to retrieve the wallet data for the authenticated user.
//
// Parameters:
//   - token: The user token used to authenticate the request.
//
// Returns:
//   - ZBDUserWalletData: The wallet data for the authenticated user.
//   - error: An error if the API request or response handling encounters issues.
//
// Example usage:
//
//	walletData, err := client.GetUserWalletData("user-token-here")
//	if err != nil {
//	  fmt.Println("Error fetching user wallet data:", err)
//	  return
//	}
//	fmt.Println("Wallet balance:", walletData.Balance)
func (c *Client) GetUserWalletData(token string) (ZBDUserWalletData, error) {
	var res ZBDUserWalletData
	headers := make(http.Header)
	headers.Set("usertoken", token)

	err := c.MakeRequest("GET", "/v1/oauth2/wallet", nil, &res, headers)
	return res, err
}

// Create Withdrawal Request: https://api-reference.zebedee.io/#60cee894-009f-40dc-9cba-e9aec5ce8aa9
//
// Takes an WithdrawalRequest object containing only
// {expiresIn, Amount, Description, InternalID, CallbackURL}
func (c *Client) WithdrawalRequest(params *WithdrawalRequest) (*WithdrawalRequest, error) {
	err := c.MakeRequest("POST", "/withdrawal-requests", params, params, nil)
	return params, err
}

// Get All Withdrawal Requests: https://api-reference.zebedee.io/#bc59c1da-4d5a-49c6-937f-f95d71c940c6
func (c *Client) ListWithdrawalRequests() ([]WithdrawalRequest, error) {
	var wr []WithdrawalRequest
	err := c.MakeRequest("GET", "/withdrawal-requests", nil, &wr, nil)
	return wr, err
}

// Get Withdrawal Request Details: https://api-reference.zebedee.io/#12aea552-0b8d-4562-a84b-a890d4f17a32
func (c *Client) GetWithdrawalRequest(wrequestID string) (*WithdrawalRequest, error) {
	var wr WithdrawalRequest
	err := c.MakeRequest("GET", "/withdrawal-requests/"+wrequestID, nil, &wr, nil)
	return &wr, err
}

// Pay Invoice: https://api-reference.zebedee.io/#04dace34-06f5-4c2f-9215-5870205098d5
//
// Takes a Payment object containing only {Description, InternalID, Invoice}
// and overwrites that with the response.
func (c *Client) Pay(params *Payment) (*Payment, error) {
	err := c.MakeRequest("POST", "/payments", params, params, nil)
	return params, err
}

// Get All Payments: https://api-reference.zebedee.io/#08ea69cc-dd6f-4381-a489-18004b911f96
func (c *Client) ListPayments() ([]Payment, error) {
	var payments []Payment
	err := c.MakeRequest("GET", "/payments", nil, &payments, nil)
	return payments, err
}

// Get Payment Details: https://api-reference.zebedee.io/#244ebe9f-6c4d-4162-a805-9a0e8955b20d
func (c *Client) GetPayment(paymentID string) (*Payment, error) {
	var payment Payment
	err := c.MakeRequest("GET", "/payments/"+paymentID, nil, &payment, nil)
	return &payment, err
}

// Send Payment to Gamertag: https://api-reference.zebedee.io/#8da3c4a3-ecf0-4fcc-be17-72e34051a1e9
func (c *Client) SendGamertagPayment(gamertag, amount, description string) (*PeerPaymentResult, error) {
	var payment PeerPaymentResult
	err := c.MakeRequest("POST", "/gamertag/send-payment", struct {
		Gamertag    string `json:"gamertag"`
		Amount      string `json:"amount"`
		Description string `json:"description"`
	}{gamertag, amount, description}, &payment, nil)
	return &payment, err
}

// Fetch Gamertag Transaction Details By ID: https://api-reference.zebedee.io/#80571b36-eac4-4966-9c49-1b83d0ae466e
func (c *Client) FetchGamerTagTransaction(transactionID string) (*PeerPayment, error) {
	var payment PeerPayment
	err := c.MakeRequest("GET", "/gamertag/transaction/"+transactionID, nil, &payment, nil)
	return &payment, err
}

// Fetch User ID By Gamertag: https://api-reference.zebedee.io/#8442d428-d4be-4082-b4a2-6e9489fe4fdf
func (c *Client) FetchUserIDFromGamertag(gamertag string) (string, error) {
	var data struct {
		ID string `json:"id"`
	}
	err := c.MakeRequest("GET", "/user-id/gamertag/"+gamertag, nil, &data, nil)
	return data.ID, err
}

// Fetch Gamertag By User ID: https://api-reference.zebedee.io/#61085d46-675f-4000-9017-9973fb1cdc80
func (c *Client) FetchGamertagFromUserID(userID string) (string, error) {
	var data struct {
		Gamertag string `json:"gamertag"`
	}
	err := c.MakeRequest("GET", "/gahttps://github.com/pocketbase/pocketbase/toolsmertag/user-id/"+userID, nil, &data, nil)
	return data.Gamertag, err
}

// Fetch Charge from Gamertag: https://api.zebedee.io/v0/gamertag/charges
func (c *Client) CreateGamertagCharge(gamertag, amount, description string) (*Charge, error) {
	var data struct {
		ID               string    `json:"id"`
		Unit             string    `json:"unit"`
		CreatedAt        time.Time `json:"createdAt"`
		Status           string    `json:"status"`
		InternalID       string    `json:"internalId"`
		Amount           string    `json:"amount"`
		Description      string    `json:"description"`
		InvoiceRequest   string    `json:"invoiceRequest"`
		InvoiceExpiresAt time.Time `json:"invoiceExpiresAt"`
	}
	err := c.MakeRequest("POST", "/gamertag/charges", struct {
		Gamertag    string `json:"gamertag"`
		Amount      string `json:"amount"`
		Description string `json:"description"`
	}{gamertag, amount, description}, &data, nil)
	if err != nil {
		return nil, err
	}

	readableStatus := data.Status
	spl := strings.Split(data.Status, "_")
	if len(spl) == 2 {
		readableStatus = strings.ToLower(spl[1])
	}

	return &Charge{
		ExpiresIn:   int64(data.InvoiceExpiresAt.Sub(time.Now()).Seconds()),
		Unit:        data.Unit,
		Amount:      data.Amount,
		Status:      readableStatus,
		Description: data.Description,
		CreatedAt:   data.CreatedAt,
		ExpiresAt:   data.InvoiceExpiresAt,
		ID:          data.ID,
		InternalID:  data.InternalID,
		Invoice: struct {
			Request string `json:"request"`
			URI     string `json:"uri"`
		}{data.InvoiceRequest, "lightning:" + data.InvoiceRequest},
	}, nil
}

// Get API Production IPs: https://api-reference.zebedee.io/#c7e18276-6935-4cca-89ae-ad949efe9a6a
func (c *Client) GetProductionIPs() ([]string, error) {
	var ips struct {
		IPs []string `json:"ips"`
	}
	err := c.MakeRequest("GET", "/prod-ips", nil, &ips, nil)
	return ips.IPs, err
}
