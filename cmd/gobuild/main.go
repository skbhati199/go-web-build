package main

import (
	"fmt"
	"os"

	"github.com/skbhati199/go-web-build/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

}
