package main

import (
	"fmt"

	"gopher.com/secret"
)

func main() {
	// Memory Vault
	v := secret.NewVault("my secret key")
	err := v.Set("my api key", "0123456789abcdef")
	if err != nil {
		panic(err)
	}

	plain, err := v.Get("my api key")
	if err != nil {
		panic(err)
	}

	fmt.Println("Value: ", plain)

	// File Vault
	fv := secret.NewVaultFile("my file secret", ".secrets")
	err = fv.SetFromFile("key", "abcdef0123456789")
	if err != nil {
		panic(err)
	}
	plain, err = fv.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println("File Value: ", plain)

}
