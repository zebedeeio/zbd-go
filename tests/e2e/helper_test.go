//go:build e2e

package e2e

import (
	"os"

	zebedee "github.com/zebedeeio/go-sdk"
)

func NewClient() *zebedee.Client {
	apiKey := os.Getenv("ZEBEDEE_API_KEY")
	return zebedee.New(apiKey)
}
