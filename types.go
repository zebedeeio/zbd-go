package zebedee

import "time"

// The Wallet Object: https://documentation.zebedee.io/docs/wallet-main/
type Wallet struct {
	Unit    string `json:"unit"`
	Balance string `json:"balance"`
}

// The Charge Object: https://documentation.zebedee.io/docs/charges-main
type Charge struct {
	ExpiresIn   time.Duration `json:"expiresIn"`
	Unit        string        `json:"unit"`
	Amount      string        `json:"amount"`
	ConfirmedAt time.Time     `json:"confirmedAt"`
	Status      string        `json:"status"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"createdAt"`
	ExpiresAt   time.Time     `json:"expiresAt"`
	ID          string        `json:"id"`
	InternalID  string        `json:"internalId"`
	CallbackURL string        `json:"callbackUrl"`
	Invoice     struct {
		Request string `json:"request"`
		URI     string `json:"uri"`
	} `json:"invoice"`
}

// The Payment Object: https://documentation.zebedee.io/docs/payments-main
type Payment struct {
	Fee         string    `json:"fee"`
	Unit        string    `json:"unit"`
	Amount      string    `json:"amount"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processedAt"`
	ID          string    `json:"id"`
	Description string    `json:"description"`
	InternalID  string    `json:"internalId"`
	Invoice     string    `json:"invoice"`
}
