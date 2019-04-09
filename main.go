package main

import (
	"fmt"
	"github.com/leanovate/mite-go/cmd"
	"os"
)

func main() {
	err := cmd.HandleCommands()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
