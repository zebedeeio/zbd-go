//go:build e2e

package e2e

import (
	"testing"

	zebedee "github.com/zebedeeio/go-sdk"
)

func TestPayments(t *testing.T) {
	client := NewClient()
	_, err := client.Pay(&zebedee.Payment{
		Description: "a payment?",
		InternalID:  "testw",
		Invoice:     "lnbc1m1p0utye7pp5xxhg0h7n6rnymjmqv09rwpplvaru83k4a6d60er04a80ts6yuc8sdqcw3jhxapq0fjkyetyv4jj6em0xq9p5hsqrzjqtqkejjy2c44jrwj08y5ygqtmn8af7vscwnflttzpsgw7tuz9r40la6l0lva5e9lvyqqqqqqqqqqqqqqpysp5qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5z5tpwxqergd3c8g7rusq9qypqsq40lkj5at0w5a7wf86hp6jr68up6u2hh9nr84ha60kaneuwr7tn2xyu6jmnjzgxkypaey4catj26q3d9lgtt0m3tc4akym4y9hp5dpcqq7ss46l",
	})
	if err == nil {
		t.Errorf(".Pay() succeeded but should have returned an error")
	}

	payments, err := client.ListPayments()
	if err != nil {
		t.Errorf("got error from .ListPayments(): %s", err)
	} else {
		if len(payments) != 1 {
			t.Errorf("should have returned one payment, but instead returned %d",
				len(payments))
		}
	}

	if len(payments) > 0 {
		payment, err := client.GetPayment(payments[0].ID)
		if err != nil {
			t.Errorf("got error from .GetPayment(): %s", err)
		} else {
			if payment.ID != payments[0].ID {
				t.Errorf("got a payment with an id different than requested")
			}
			if payment.Amount != payments[0].Amount ||
				payment.Fee != payments[0].Fee ||
				payment.Status != payments[0].Status ||
				payment.Invoice != payments[0].Invoice ||
				payment.Description != payments[0].Description ||
				payment.InternalID != payments[0].InternalID {
				t.Errorf("payment returned from .ListPayments() is different from the" +
					" same payment returned from .GetPayment()")
			}
			if payment.Amount != "10000" || payment.Description != "test zebedee-go" {
				t.Errorf("details of the payment are wrong")
			}
		}
	}
}

func TestPaymentsBad(t *testing.T) {
	client := NewClient()
	_, err := client.Pay(&zebedee.Payment{
		Invoice: "x",
	})
	if err == nil {
		t.Errorf(".Pay() should have returned an error")
	}
}
