package main

import (
	"fmt"
	"os"

	"github.com/apus3404-oss/xupav-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
