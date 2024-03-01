//go:build e2e

package e2e

import zebedee "github.com/zebedeeio/go-sdk"

func NewClient() *zebedee.Client {
	return zebedee.New("edg7SOTFWbh1FbjVecbmZi4G4nYVHJj2")
}
