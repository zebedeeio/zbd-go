package main

import (
	"fmt"
	"time"

	zebedee "github.com/zebedeeio/zebedee-go"
)

func main() {
	zbd := zebedee.New("")

	// generate a charge
	charge, _ := zbd.Charge(&zebedee.Charge{
		InternalID:  "aninternalchargeid",
		Amount:      "10000",
		Description: "My charge description",
		ExpiresIn:   time.Minute * 30,
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
