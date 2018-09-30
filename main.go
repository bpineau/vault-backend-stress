package main

import (
	"fmt"
	"os"

	"github.com/bpineau/vault-backend-stress/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}
