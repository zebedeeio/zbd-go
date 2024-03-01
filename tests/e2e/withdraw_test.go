//go:build e2e

package e2e

import (
	"sort"
	"strings"
	"testing"
	"time"

	zebedee "github.com/zebedeeio/go-sdk"
)

func TestWithdrawalRequests(t *testing.T) {
	client := NewClient()
	wr, err := client.WithdrawalRequest(&zebedee.WithdrawalRequest{
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
	client := NewClient()
	_, err := client.WithdrawalRequest(&zebedee.WithdrawalRequest{
		Amount:      "5000000",
		Description: "a test withdrawal request",
		InternalID:  "testd",
		CallbackURL: "https://example.com/callback",
	})
	if err == nil {
		t.Errorf(".WithdrawalRequest() should have returned an error")
	}
}
