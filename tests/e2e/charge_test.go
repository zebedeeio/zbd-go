//go:build e2e

package e2e

import (
	"sort"
	"strings"
	"testing"
	"time"

	zebedee "github.com/zebedeeio/go-sdk"
)

func TestCharges(t *testing.T) {
	client := NewClient()
	charge, err := client.Charge(&zebedee.Charge{
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
	client := NewClient()
	_, err := client.Charge(&zebedee.Charge{
		Amount:      "123000000",
		Description: "a test invoice",
		InternalID:  "testb",
	})
	if err == nil {
		t.Errorf(".Charge() should have returned an error")
	} else if err.Error() != "The maximum Charge amount supported is 45,000 satoshis." {
		t.Errorf(".Charge() returned the wrong error")
	}

	_, err = client.Charge(&zebedee.Charge{
		Amount:      "-120",
		Description: "a test invoice",
		InternalID:  "testb",
	})
	if err == nil {
		t.Errorf(".Charge() should have returned an error")
	}
}
