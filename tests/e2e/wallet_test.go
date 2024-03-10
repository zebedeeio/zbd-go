//go:build e2e

package e2e

import (
	"testing"

	zebedee "github.com/zebedeeio/go-sdk"
)

func TestWallet(t *testing.T) {
	client := NewClient()
	_, err := client.Wallet()
	if err != nil {
		t.Errorf("got error from .Wallet(): %s", err)
		return
	}
}

func TestBadAuth(t *testing.T) {
	badClient := zebedee.New("invalidkey")
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
