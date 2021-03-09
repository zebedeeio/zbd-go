# go-sdk

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/zebedeeio/go-sdk) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](LICENSE)

ZEBEDEE's `go-sdk` is a Go package that provides a simple wrapper over the [ZEBEDEE API](https://documentation.zebedee.io/).

## Installing

```
go get github.com/zebedeeio/go-sdk
```

## Usage

Below is an example which shows some common use cases for `zebedee`. More documentation can be found on [Godoc](https://godoc.org/github.com/zebedeeio/go-sdk).

```golang
package main

import zebedee "github.com/zebedeeio/go-sdk"

func main () {
	zbd := zebedee.New(APIKEY)

	// generate a charge
	charge, _ := zbd.Charge(&zebedee.Charge{
		InternalID:  "aninternalchargeid",
		Amount:      "10000",
		Description: "My charge description",
		ExpiresIn:   int64((time.Minute * 30).Seconds()),
		CallbackURL: "https://mysite.example.com/callback/zebedee",
	})
	fmt.Printf("invoice to be paid: %s\n", charge.Invoice.Request)

	// fetch charge details
	charge, _ = zbd.GetCharge(charge.ID)
	if charge.Status == "pending" {
		fmt.Printf("payment still pending. expiring in %d seconds\n",
			int64(charge.ExpiresAt.Sub(time.Now()).Seconds()))
	}

	// pay an invoice
	payment, _ := zbd.Pay(&zebedee.Payment{
		InternalID:  "aninternalpaymentid",
		Description: "a payment I'm making",
		Invoice:     "lnbc...",
	})
	fmt.Printf("payment sent: %s\n", payment.Status)

	// get wallet balance
	wallet, _ := zbd.Wallet()
	fmt.Printf("your balance is %s %s\n", wallet.Balance, wallet.Unit)
}
```
