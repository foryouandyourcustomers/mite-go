package main

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/leanovate/mite-go/mite"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	servicesCommand.AddCommand(listServicesCommand)
	rootCmd.AddCommand(servicesCommand)
}

var servicesCommand = &cobra.Command{
	Use:   "services",
	Short: "list & adds services",
	Run:   listServicesCommand.Run,
}

var listServicesCommand = &cobra.Command{
	Use:   "list",
	Short: "list services",
	Run: func(cmd *cobra.Command, args []string) {
		api := mite.NewMiteApi(configGetApiUrl(), configGetApiKey())
		services, err := api.Services()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return
		}

		t := tabby.New()
		t.AddHeader("id", "name", "notes")
		for _, service := range services {
			t.AddLine(service.Id, service.Name, service.Note)
		}
		t.Print()
	},
}
