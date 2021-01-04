package main

import (
	"fmt"
	"os"

	"github.com/tumani1/diexample/di/sarulabsdingo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
