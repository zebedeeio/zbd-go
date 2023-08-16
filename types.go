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

type SendLightningAddressPaymentOptionsType struct {
	LnAddress   string `json:"lnAddress"`
	Amount      string `json:"amount"`
	Comment     string `json:"comment"`
	CallbackUrl string `json:"callbackUrl"`
	InternalID  string `json:"internalId"`
}

type SendLightningAddressPaymentDataResponseType struct {
	Data struct {
		ID            string `json:"id"`
		Fee           string `json:"fee"`
		Unit          string `json:"unit"`
		Amount        string `json:"amount"`
		Invoice       string `json:"invoice"`
		Preimage      string `json:"preimage"`
		WalletID      string `json:"walletId"`
		TransactionID string `json:"transactionId"`
		CallbackUrl   string `json:"callbackUrl"`
		InternalID    string `json:"internalId"`
		Comment       string `json:"comment"`
		ProcessedAt   string `json:"processedAt"`
		CreatedAt     string `json:"createdAt"`
		Status        string `json:"status"`
	} `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ValidateLightningAddressDataResponseType struct {
	Data struct {
		Valid    bool `json:"valid"`
		Metadata struct {
			MinSendable    int    `json:"minSendable"`
			MaxSendable    int    `json:"maxSendable"`
			CommentAllowed int    `json:"commentAllowed"`
			Tag            string `json:"tag"`
			Metadata       string `json:"metadata"`
			Callback       string `json:"callback"`
			PayerData      struct {
				Name struct {
					Mandatory bool `json:"mandatory"`
				} `json:"name"`
				Identifier struct {
					Mandatory bool `json:"mandatory"`
				} `json:"identifier"`
			} `json:"payerData"`
			Disposable bool `json:"disposable"`
		} `json:"metadata"`
	} `json:"data"`
	Success bool `json:"success"`
}

type CreateChargeFromLightningAddressOptionsType struct {
	Amount      string `json:"amount"`
	LNAddress   string `json:"lnaddress"`
	Description string `json:"description"`
}

type FetchChargeFromLightningAddressDataResponseType struct {
	Data struct {
		LNAddress string `json:"lnaddress"`
		Amount    string `json:"amount"`
		Invoice   struct {
			URI     string `json:"uri"`
			Request string `json:"request"`
		} `json:"invoice"`
	} `json:"data"`
	Success bool `json:"success"`
}

type KeysendDataResponseType struct {
	Data struct {
		KeysendID   string `json:"keysendId"`
		PaymentID   string `json:"paymentId"`
		Transaction struct {
			ID          string `json:"id"`
			WalletID    string `json:"walletId"`
			Type        string `json:"type"`
			TotalAmount string `json:"totalAmount"`
			Fee         string `json:"fee"`
			Amount      string `json:"amount"`
			Description string `json:"description"`
			Status      string `json:"status"`
			ConfirmedAt string `json:"confirmedAt"`
		} `json:"transaction"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type KeysendOptionsType struct {
	Amount      string `json:"amount"`
	Pubkey      string `json:"pubkey"`
	TLVRecords  string `json:"tlvRecords"`
	Metadata    string `json:"metadata"`
	CallbackURL string `json:"callbackUrl"`
}

type BTCUSDDataResponseType struct {
	Data struct {
		BTCUSDPrice     string `json:"btcUsdPrice"`
		BTCUSDTimestamp string `json:"btcUsdTimestamp"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type SupportedRegionDataResponseType struct {
	Data struct {
		IPAddress   string `json:"ipAddress"`
		IsSupported bool   `json:"isSupported"`
		IPCountry   string `json:"ipCountry"`
		IPRegion    string `json:"ipRegion"`
	} `json:"data"`
	Success bool `json:"success"`
}

type ProdIPSDataResponseType struct {
	Data struct {
		IPS []string `json:"ips"`
	} `json:"data"`
	Success bool `json:"success"`
}
