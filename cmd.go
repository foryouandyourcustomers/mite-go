package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mite-go",
	Short: "cli client for mite time tracking",
	Run: func(cmd *cobra.Command, args []string) {
		// list entries for last 7 days
	},
}

func cmdLineHandler() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
