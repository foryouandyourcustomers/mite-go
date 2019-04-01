package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	listCommand.AddCommand(listProjectsCommand)
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "list entries, projects and roles",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: list entries for last 7 days by default
	},
}

var listProjectsCommand = &cobra.Command{
	Use:   "projects",
	Short: "list projects",
	Run: func(cmd *cobra.Command, args []string) {
		url := configGetApiUrl()
		key := configGetApiKey()
		// do get request
		fmt.Println(url, key)
	},
}
