package zebedee

import (
	"sort"
	"strings"
	"testing"
	"time"
)

var client *Client

func TestMain(m *testing.M) {
	client = New("edg7SOTFWbh1FbjVecbmZi4G4nYVHJj2", nil)
	client.BaseURL = "https://dev.zebedee.io/v0"
	m.Run()
}

func TestWallet(t *testing.T) {
	_, err := client.Wallet()
	if err != nil {
		t.Errorf("got error from .Wallet(): %s", err)
		return
	}
}

func TestBadAuth(t *testing.T) {
	badClient := New("invalidkey", nil)
	badClient.BaseURL = "https://dev.zebedee.io/v0"

	_, err := badClient.Wallet()
	if err == nil {
		t.Errorf("should have gotten an error from .Wallet()")
	}

	const errorMessage = "Invalid authentication credentials"
	if err.Error() != errorMessage {
		t.Errorf("error was '%s', should have been '%s'", err.Error(), errorMessage)
	}
}

func TestCharges(t *testing.T) {
	charge, err := client.Charge(&Charge{
		ExpiresIn:   30 * 60,
		Amount:      "123000",
		Description: "a test invoice",
		InternalID:  "testx",
		CallbackURL: "https://example.com/callback",
	})
	if err != nil {
		t.Errorf("got error from .Charge(): %s", err)
	} else {
		if charge.ExpiresAt.After(time.Now().Add(time.Minute * 30)) {
			t.Error("charge expires after we wanted")
		}
		if !strings.HasPrefix(charge.Invoice.Request, "lnbc123") {
			t.Error("charge has wrong bolt11 invoice")
		}
	}

	// fetch this same charge
	charge, err = client.GetCharge(charge.ID)
	if err != nil {
		t.Errorf("got error from .GetCharge(): %s", err)
	} else {
		if charge.Amount != "123000" {
			t.Error("charge amount is different than specified")
		}
		if charge.Description != "a test invoice" {
			t.Error("charge description is different than specified")
		}
		if charge.InternalID != "testx" {
			t.Error("charge internal id is different than specified")
		}
		if charge.CallbackURL != "https://example.com/callback" {
			t.Error("charge callback url is different than specified")
		}
	}

	// get all charges
	charges, err := client.ListCharges()
	if err != nil {
		t.Errorf("got error from .ListCharges(): %s", err)
	} else {
		sort.Slice(charges, func(i, j int) bool { return charges[i].CreatedAt.Before(charges[j].CreatedAt) })
		if charges[len(charges)-1].ID != charge.ID {
			t.Errorf("last charge from list is not the charge we just created (%s != %s)",
				charges[len(charges)-1].ID, charge.ID)
		}
	}
}

func TestChargesBad(t *testing.T) {
	_, err := client.Charge(&Charge{
		Amount:      "123000000",
		Description: "a test invoice",
		InternalID:  "testb",
	})
	if err == nil {
		t.Errorf(".Charge() should have returned an error")
	} else if err.Error() != "The maximum Charge amount supported is 45,000 satoshis." {
		t.Errorf(".Charge() returned the wrong error")
	}

	_, err = client.Charge(&Charge{
		Amount:      "-120",
		Description: "a test invoice",
		InternalID:  "testb",
	})
	if err == nil {
		t.Errorf(".Charge() should have returned an error")
	}
}

func TestWithdrawalRequests(t *testing.T) {
	wr, err := client.WithdrawalRequest(&WithdrawalRequest{
		ExpiresIn:   30 * 60,
		Amount:      "50000",
		Description: "a test withdrawal request",
		InternalID:  "testy",
		CallbackURL: "https://example.com/callback",
	})
	if err != nil {
		t.Errorf("got error from .WithdrawalRequest(): %s", err)
	} else {
		if wr.ExpiresAt.After(time.Now().Add(time.Minute * 30)) {
			t.Error("wr expires after we wanted")
		}
		if !strings.HasPrefix(wr.Invoice.Request, "lnurl1") {
			t.Errorf("wr has something that isn't an lnurl: '%s'", wr.Invoice.Request)
		}
		if !strings.HasPrefix(wr.Invoice.FastRequest, "lnurl1") {
			t.Errorf("wr has something that isn't a fast lnurl: '%s'", wr.Invoice.Request)
		}
	}

	// fetch this same wr
	wr, err = client.GetWithdrawalRequest(wr.ID)
	if err != nil {
		t.Errorf("got error from .GetWithdrawalRequest(): %s", err)
	} else {
		if wr.Status != "pending" {
			t.Errorf("wr is not pending")
		}
		if wr.Amount != "50000" {
			t.Error("wr amount is different than specified")
		}
		if wr.Description != "a test withdrawal request" {
			t.Error("wr description is different than specified")
		}
		if wr.InternalID != "testy" {
			t.Error("wr internal id is different than specified")
		}
		if wr.CallbackURL != "https://example.com/callback" {
			t.Error("wr callback url is different than specified")
		}
	}

	// get all wrs
	wrs, err := client.ListWithdrawalRequests()
	if err != nil {
		t.Errorf("got error from .ListWithdrawalRequests(): %s", err)
	} else {
		sort.Slice(wrs, func(i, j int) bool { return wrs[i].CreatedAt.Before(wrs[j].CreatedAt) })
		if wrs[len(wrs)-1].ID != wr.ID {
			t.Errorf("last wr from list is not the wr we just created (%s != %s)",
				wrs[len(wrs)-1].ID, wr.ID)
		}
	}
}

func TestWithdrawalRequestsBad(t *testing.T) {
	_, err := client.WithdrawalRequest(&WithdrawalRequest{
		Amount:      "5000000",
		Description: "a test withdrawal request",
		InternalID:  "testd",
		CallbackURL: "https://example.com/callback",
	})
	if err == nil {
		t.Errorf(".WithdrawalRequest() should have returned an error")
	}
}

func TestPayments(t *testing.T) {
	_, err := client.Pay(&Payment{
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
	_, err := client.Pay(&Payment{
		Invoice: "x",
	})
	if err == nil {
		t.Errorf(".Pay() should have returned an error")
	}
}

func TestDecodeCharge(t *testing.T) {
	response, err := client.DecodeCharge(&DecodeChargeOptionsType{
		Invoice: "An invoice",
	})

	if err != nil {
		t.Errorf("got error from .DecodeCharge(): %v", err)
	}

	if !response.Success {
		t.Errorf("unexpected success value: %v", response.Success)
	}
	if response.Data.Unit != "BTC" {
		t.Errorf("unexpected unit value: %s", response.Data.Unit)
	}
}
