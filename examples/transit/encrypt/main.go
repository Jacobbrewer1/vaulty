package main

import (
	"context"
	"fmt"

	"github.com/Jacobbrewer1/vaulty"
)

const (
	vaultAddr           = "http://localhost:8200"
	vaultPrefix         = "some-prefix" // No need to include the path type here (e.g. just "some-prefix" instead of "transit/encrypt")
	vaultTransitKeyName = "key-name"
	vaultDecryptedData  = "some-data" // Data that has been programmatically read from the DB, etc.

	// Read these from a config file in production
	vaultUser = "username"
	vaultPass = "password"
)

func main() {
	vc, err := vaulty.NewClient(
		vaulty.WithGeneratedVaultClient(vaultAddr),
		vaulty.WithUserPassAuth(vaultUser, vaultPass),
	)
	if err != nil {
		panic(err)
	}

	sec, err := vc.Path(vaultTransitKeyName, vaulty.WithPrefix(vaultPrefix)).TransitEncrypt(context.Background(), vaultDecryptedData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Encrypted data:", vaulty.GetTransitCipherText(sec))
}
