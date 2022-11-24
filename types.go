package zebedee

import "time"

type Wallet struct {
	Unit    string `json:"unit"`
	Balance string `json:"balance"`
}

type Charge struct {
	ExpiresIn       int64     `json:"expiresIn"`
	Unit            string    `json:"unit"`
	Amount          string    `json:"amount"`
	ConfirmedAt     time.Time `json:"confirmedAt"`
	Status          string    `json:"status"`
	Description     string    `json:"description"`
	DescriptionHash string    `json:"invoiceDescriptionHash"`
	CreatedAt       time.Time `json:"createdAt"`
	ExpiresAt       time.Time `json:"expiresAt"`
	ID              string    `json:"id"`
	InternalID      string    `json:"internalId"`
	CallbackURL     string    `json:"callbackUrl"`
	Invoice         struct {
		Request string `json:"request"`
		URI     string `json:"uri"`
	} `json:"invoice"`
}

type WithdrawalRequest struct {
	ExpiresIn   int64     `json:"expiresIn"`
	Unit        string    `json:"unit"`
	Amount      string    `json:"amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	Description string    `json:"description"`
	ID          string    `json:"id"`
	InternalID  string    `json:"internalId"`
	CallbackURL string    `json:"callbackUrl"`
	Invoice     struct {
		Request     string `json:"request"`
		FastRequest string `json:"fastRequest"`
		URI         string `json:"uri"`
		FastURI     string `json:"fastUri"`
	} `json:"invoice"`
}

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

type PeerPaymentResult struct {
	ReceiverId    string `json:"receiverId"`
	TransactionId string `json:"transactionId"`
	Amount        string `json:"amount"`
	Comment       string `json:"comment"`
}

type PeerPayment struct {
	ID          string    `json:"id"`
	ReceiverID  string    `json:"receiverId"`
	Amount      string    `json:"amount"`
	Fee         string    `json:"fee"`
	Unit        string    `json:"unit"`
	ProcessedAt time.Time `json:"processedAt"`
	ConfirmedAt time.Time `json:"confirmedAt"`
	Comment     string    `json:"comment"`
	Status      string    `json:"status"`
}
