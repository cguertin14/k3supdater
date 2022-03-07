package main

import (
	"fmt"
	"os"

	"github.com/cguertin14/k3supdater/cmd"
)

func main() {
	// Execute program and if there is an error,
	// show it to the user.
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
