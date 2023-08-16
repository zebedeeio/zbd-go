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

type DecodeChargeOptionsType struct {
	Invoice string `json:"invoice"`
}

type DecodeChargeResponseType struct {
	Data struct {
		Unit                   string `json:"unit"`
		Status                 string `json:"status"`
		Amount                 string `json:"amount"`
		CreatedAt              string `json:"createdAt"`
		InternalId             string `json:"internalId"`
		CallbackUrl            string `json:"callbackUrl"`
		Description            string `json:"description"`
		InvoiceRequest         string `json:"invoiceRequest"`
		InvoiceExpiresAt       string `json:"invoiceExpiresAt"`
		InvoiceDescriptionHash string `json:"invoiceDescriptionHash,omitempty"`
	} `json:"data"`
	Success bool `json:"success"`
}

type StaticChargeOptionsType struct {
	AllowedSlots   *string `json:"allowedSlots"`
	MinAmount      string  `json:"minAmount"`
	MaxAmount      string  `json:"maxAmount"`
	Description    string  `json:"description"`
	InternalID     string  `json:"internalId"`
	CallbackURL    string  `json:"callbackUrl"`
	SuccessMessage string  `json:"successMessage"`
}

type StaticChargeDataResponseType struct {
	Data struct {
		ID             string  `json:"id"`
		Unit           string  `json:"unit"`
		Slots          string  `json:"slots"`
		MinAmount      string  `json:"minAmount"`
		MaxAmount      string  `json:"maxAmount"`
		CreatedAt      string  `json:"createdAt"`
		CallbackURL    string  `json:"callbackUrl"`
		InternalID     string  `json:"internalId"`
		Description    string  `json:"description"`
		ExpiresAt      string  `json:"expiresAt"`
		ConfirmedAt    string  `json:"confirmedAt"`
		SuccessMessage string  `json:"successMessage"`
		AllowedSlots   *string `json:"allowedSlots"`
		Status         string  `json:"status"`
		Fee            string  `json:"fee"`
		Invoice        struct {
			Request string `json:"request"`
			URI     string `json:"uri"`
		} `json:"invoice"`
	} `json:"data"`
	Message string `json:"message"`
}
