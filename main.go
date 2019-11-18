package main

import (
	"fmt"
	"os"

	"gitlab.com/igor.tumanov1/theboatscom/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
