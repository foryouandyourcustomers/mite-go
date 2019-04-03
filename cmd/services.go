package cmd

import (
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

func init() {
	servicesCommand.AddCommand(listServicesCommand)
	rootCmd.AddCommand(servicesCommand)
}

var servicesCommand = &cobra.Command{
	Use:   "services",
	Short: "list & adds services",
	RunE:  listServicesCommand.RunE,
}

var listServicesCommand = &cobra.Command{
	Use:   "list",
	Short: "list services",
	RunE: func(cmd *cobra.Command, args []string) error {
		services, err := deps.miteApi.Services()
		if err != nil {
			return err
		}

		t := tabby.New()
		t.AddHeader("id", "name", "notes")
		for _, service := range services {
			t.AddLine(service.Id, service.Name, service.Note)
		}
		t.Print()
		return nil
	},
}
