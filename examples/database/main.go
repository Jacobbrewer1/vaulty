package main

import (
	"context"
	"fmt"

	"github.com/jacobbrewer1/vaulty"
)

const (
	vaultAddr     = "http://localhost:8200"
	vaultDBPath   = "role-name"
	vaultDBPrefix = "mount/creds" // We need the path type here (e.g. "mount/creds" with the "creds" path type)

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

	sec, err := vc.Path(vaultDBPath, vaulty.WithPrefix(vaultDBPrefix)).GetSecret(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(sec.Data["username"])
	fmt.Println(sec.Data["password"])
}
