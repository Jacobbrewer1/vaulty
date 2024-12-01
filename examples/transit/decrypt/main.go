package main

import (
	"context"
	"fmt"

	"github.com/jacobbrewer1/vaulty"
)

const (
	vaultAddr           = "http://localhost:8200"
	vaultPrefix         = "some-prefix" // No need to include the path type here (e.g. just "some-prefix" instead of "transit/decrypt")
	vaultTransitKeyName = "key-name"
	vaultEncryptedData  = "some-data" // Data that has been programmatically read from the DB, etc.

	// Read these from a config file in production
	vaultUser = "username"
	vaultPass = "password"
)

func main() {
	vc, err := vaulty.NewClient(
		vaulty.WithAddr(vaultAddr),
		vaulty.WithUserPassAuth(vaultUser, vaultPass),
	)
	if err != nil {
		panic(err)
	}

	sec, err := vc.Path(vaultTransitKeyName, vaulty.WithPrefix(vaultPrefix)).TransitDecrypt(context.Background(), vaultEncryptedData)
	if err != nil {
		panic(err)
	}

	fmt.Println(sec)
}
