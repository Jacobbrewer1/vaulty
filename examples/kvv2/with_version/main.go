package main

import (
	"context"
	"fmt"

	"github.com/jacobbrewer1/vaulty"
)

const (
	vaultAddr      = "http://localhost:8200"
	vaultKVName    = "secret-name"
	vaultKVV2Mount = "secret-mount"
	vaultKVVersion = 0 // 0 is the latest version; specify a version number to get a specific version

	// Read these from a config file in production
	vaultUser = "username"
	vaultPass = "password"
)

func main() {
	vc, err := vaulty.NewClient(
		vaulty.WithAddr(vaultAddr),
		vaulty.WithUserPassAuth(vaultUser, vaultPass),
		vaulty.WithKvv2Mount(vaultKVV2Mount),
	)
	if err != nil {
		panic(err)
	}

	sec, err := vc.Path(vaultKVName, vaulty.WithVersion(vaultKVVersion)).GetKvSecretV2(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(sec.Data)
}
