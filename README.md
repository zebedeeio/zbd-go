## Usage

```golang
package main

import zebedee "github.com/zebedeeio/zebedee-go"

func main () {
    zbd := zebedee.New(APIKEY)

    // generate a charge
    charge, _ := zbd.Charge(&zebedee.Charge{
        InternalId: "aninternalchargeid",
        Amount: 10000,
        Description: "My charge description",
        ExpiresIn: time.Minute * 30,
        CallbackURL: "https://mysite.example.com/callback/zebedee",
    })
    fmt.Printf("invoice to be paid: %s\n", charge.Invoice.Request)

    // fetch charge details
    charge, _ := zbd.GetCharge(charge.ID)
    if charge.Status == "pending" {
        fmt.Printf("payment still pending. expiring in %d seconds\n",
            charge.ExpiresAt.Sub(time.Now()).Seconds())
    }

    // pay an invoice
    payment, err := zbd.Pay(&zebedee.Payment{
        InternalId: "aninternalpaymentid",
        Description: "a payment I'm making",
        Invoice: "lnbc...",
    })
    fmt.Printf("payment sent: %s\n", payment.Status)
}
```
