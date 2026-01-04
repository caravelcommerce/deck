package main

import (
	"fmt"
	"os"

	"github.com/caravelcommerce/deck/cmd"
)

// Version é definida em tempo de build
var Version = "dev"

func main() {
	cmd.SetVersion(Version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
